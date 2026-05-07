package io.github.metamessage.mm;

public final class CamelToSnake {
    private CamelToSnake() {}

    public static String convert(String s) {
        if (s == null || s.isEmpty()) {
            return "";
        }
        StringBuilder result = new StringBuilder(s.length() + 4);
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            if (Character.isUpperCase(c)) {
                if (i > 0) {
                    char prev = s.charAt(i - 1);
                    boolean prevUpper = Character.isUpperCase(prev);
                    boolean nextUpper = i + 1 < s.length() && Character.isUpperCase(s.charAt(i + 1));
                    if (!prevUpper || (i + 1 < s.length() && !nextUpper)) {
                        result.append('_');
                    }
                }
                result.append(Character.toLowerCase(c));
            } else {
                result.append(c);
            }
        }
        return result.toString();
    }
}
