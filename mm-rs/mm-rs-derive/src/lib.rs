use proc_macro::TokenStream;
use proc_macro2::{TokenStream as TokenStream2, TokenTree};
use quote::quote;
use syn::{parse_macro_input, Data, DeriveInput, Fields};

fn extract_mm_attrs(attrs: &[syn::Attribute]) -> Vec<(String, String)> {
    let mut pairs = Vec::new();
    for attr in attrs {
        if !attr.path().is_ident("mm") {
            continue;
        }
        if let syn::Meta::List(meta_list) = &attr.meta {
            pairs = parse_mm_pairs(&meta_list.tokens);
        }
    }
    pairs
}

fn parse_mm_pairs(tokens: &TokenStream2) -> Vec<(String, String)> {
    let mut pairs = Vec::new();
    let mut tokens_iter = tokens.clone().into_iter().peekable();

    while tokens_iter.peek().is_some() {
        let key = match tokens_iter.next() {
            Some(TokenTree::Ident(ident)) => ident.to_string(),
            Some(TokenTree::Punct(punct)) if punct.as_char() == ',' => {
                continue;
            }
            _ => {
                break;
            }
        };

        if let Some(TokenTree::Punct(p)) = tokens_iter.peek() {
            if p.as_char() == '=' {
                tokens_iter.next();
                let value = match tokens_iter.next() {
                    Some(TokenTree::Literal(lit)) => {
                        let s = lit.to_string();
                        if s.starts_with('"') && s.ends_with('"') {
                            s[1..s.len() - 1].to_string()
                        } else {
                            s
                        }
                    }
                    Some(TokenTree::Ident(ident)) => ident.to_string(),
                    other => {
                        if let Some(ref t) = other {
                            format!("{:?}", t)
                        } else {
                            String::new()
                        }
                    }
                };
                pairs.push((key, value));
            } else {
                pairs.push((key, String::new()));
            }
        } else {
            pairs.push((key, String::new()));
        }
    }

    pairs
}

fn build_tag_code(pairs: &[(String, String)]) -> proc_macro2::TokenStream {
    let mut setters = Vec::new();

    for (key, value) in pairs {
        match key.as_str() {
            "name" => {
                let v = value.clone();
                setters.push(quote! { tag.name = Some(#v.to_string()); });
            }
            "desc" => {
                let v = value.clone();
                setters.push(quote! { tag.desc = Some(#v.to_string()); });
            }
            "type" => {
                let v = value.clone();
                setters.push(quote! { tag.value_type = crate::ir::ValueType::from_str(#v); });
            }
            "nullable" => {
                if value.is_empty() || value == "true" {
                    setters.push(quote! { tag.nullable = true; });
                } else {
                    setters.push(quote! { tag.nullable = false; });
                }
            }
            "allow_empty" => {
                if value.is_empty() || value == "true" {
                    setters.push(quote! { tag.allow_empty = true; });
                } else {
                    setters.push(quote! { tag.allow_empty = false; });
                }
            }
            "unique" => {
                if value.is_empty() || value == "true" {
                    setters.push(quote! { tag.unique = true; });
                } else {
                    setters.push(quote! { tag.unique = false; });
                }
            }
            "deprecated" => {
                if value.is_empty() || value == "true" {
                    setters.push(quote! { tag.deprecated = true; });
                } else {
                    setters.push(quote! { tag.deprecated = false; });
                }
            }
            "default_val" => {
                let v = value.clone();
                setters.push(quote! { tag.default_val = Some(#v.to_string()); });
            }
            "min" => {
                let v = value.clone();
                setters.push(quote! { tag.min = Some(#v.to_string()); });
            }
            "max" => {
                let v = value.clone();
                setters.push(quote! { tag.max = Some(#v.to_string()); });
            }
            "size" => {
                let v = value.clone();
                if let Ok(n) = v.parse::<u64>() {
                    setters.push(quote! { tag.size = Some(#n); });
                }
            }
            "enums" => {
                let v = value.clone();
                setters.push(quote! { tag.enums = Some(#v.to_string()); });
                setters.push(quote! { tag.value_type = crate::ir::ValueType::Enum; });
            }
            "pattern" => {
                let v = value.clone();
                setters.push(quote! { tag.pattern = Some(#v.to_string()); });
            }
            "location" => {
                let v = value.clone();
                if let Ok(n) = v.parse::<i32>() {
                    setters.push(quote! { tag.location = Some(#n); });
                }
            }
            "version" => {
                let v = value.clone();
                if let Ok(n) = v.parse::<i32>() {
                    setters.push(quote! { tag.version = Some(#n); });
                }
            }
            "mime" => {
                let v = value.clone();
                setters.push(quote! { tag.mime = Some(#v.to_string()); });
            }
            _ => {}
        }
    }

    quote! {
        let mut tag = crate::ir::Tag::new();
        #(#setters)*
    }
}

fn camel_to_snake(name: &str) -> String {
    let mut result = String::new();
    for (i, c) in name.chars().enumerate() {
        if c.is_uppercase() {
            if i > 0 {
                result.push('_');
            }
            result.push(c.to_ascii_lowercase());
        } else {
            result.push(c);
        }
    }
    result
}

#[proc_macro_derive(ToNode, attributes(mm))]
pub fn derive_to_node(input: TokenStream) -> TokenStream {
    let input = parse_macro_input!(input as DeriveInput);
    let struct_name = &input.ident;

    let struct_mm_attrs = extract_mm_attrs(&input.attrs);
    let struct_tag_code = build_tag_code(&struct_mm_attrs);

    let fields = match &input.data {
        Data::Struct(data) => match &data.fields {
            Fields::Named(fields) => &fields.named,
            _ => {
                return quote! {
                    compile_error!("ToNode can only be derived for structs with named fields");
                }
                .into();
            }
        },
        _ => {
            return quote! {
                compile_error!("ToNode can only be derived for structs");
            }
            .into();
        }
    };

    let mut field_conversions = Vec::new();

    for field in fields {
        let field_name = field.ident.as_ref().unwrap();
        let field_key = camel_to_snake(&field_name.to_string());
        let field_mm_attrs = extract_mm_attrs(&field.attrs);
        let field_tag_code = build_tag_code(&field_mm_attrs);

        let conversion = quote! {
            {
                #field_tag_code
                tag.name = Some(#field_key.to_string());
                let field_node = crate::core::value_to_node::value_to_node(
                    &self.#field_name,
                    Some(tag),
                );
                crate::ir::Field {
                    key: #field_key.to_string(),
                    value: field_node,
                }
            }
        };
        field_conversions.push(conversion);
    }

    let struct_name_snake = camel_to_snake(&struct_name.to_string());

    let expanded = quote! {
        impl crate::core::value_to_node::ToNode for #struct_name {
            fn to_node(&self, top_tag: Option<crate::ir::Tag>) -> crate::ir::Node {
                use crate::core::value_to_node::ToNode;

                let mut tag = crate::ir::Tag::new();
                #struct_tag_code;

                if let Some(t) = top_tag {
                    tag = crate::ir::Tag::merge(Some(tag), t);
                }

                tag.value_type = crate::ir::ValueType::Obj;
                tag.name = Some(#struct_name_snake.to_string());
                let path = #struct_name_snake.to_string();

                let fields = vec![
                    #(#field_conversions),*
                ];

                crate::ir::Node::Object(crate::ir::Object {
                    fields,
                    tag: Some(tag),
                    path,
                })
            }
        }
    };

    TokenStream::from(expanded)
}
