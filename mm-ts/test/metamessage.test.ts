import { ValueType } from '../src/ast/value-type.js';
import { decodeToValue, encodeFromValue } from '../src/metamessage.js';
import { mm, toMM } from '../src/mm/mm.js';
import { ValueToNode } from '../src/mm/value-to-node.js';

// npm test test/metamessage.test.ts -- -t "encode and decode"
describe('Build Value', () => {
  describe('obj', () => {
    test('', () => {
      const obj = mm.obj({
        '2323': mm.str('hello'),
      });
      const obj1 = mm.obj({
        1212: mm.str('hello'),
      });
      const obj3 = mm.obj({
        '2323': mm.obj({
          '2323': mm.str('hello'),
        }),
      });
    });
  });

  describe('map', () => {
    test('', () => {
      const obj = mm.map({
        '2323': mm.str('hello'),
      });
      const obj3 = mm.map({
        '2323': mm.map({
          '2323': mm.str('hello'),
        }),
      });
    });
  });

  describe('slice', () => {
    test('', () => {
      const obj = mm.slice([mm.str('hello')]);
      const obj3 = mm.slice([mm.slice([mm.str('hello')])]);
    });
  });

  describe('arr', () => {
    test('', () => {
      const obj = mm.arr([mm.str('hello')]);
      const obj3 = mm.arr([mm.slice([mm.str('hello')])]);
    });
  });

  describe('str', () => {
    test('', () => {
      const v = mm.str('hello');
      const v2 = mm.str('');
      const v3 = mm.str('test string with special chars: @#$%^&*()');
    });
  });

  describe('bool', () => {
    test('', () => {
      const v = mm.bool(true);
      const v2 = mm.bool(false);
    });
  });

  describe('bytes', () => {
    test('', () => {
      const v2 = mm.bytes(new Uint8Array([0x01, 0x02, 0x03]));
    });
  });

  describe('int', () => {
    test('', () => {
      const v = mm.i(42n);
      const v2 = mm.i(0n);
      const v3 = mm.i(-1n);
    });
  });

  describe('int8', () => {
    test('', () => {
      const v = mm.i8(127);
      const v2 = mm.i8(-128);
    });
  });

  describe('int16', () => {
    test('', () => {
      const v = mm.i16(32767);
      const v2 = mm.i16(-32768);
    });
  });

  describe('int32', () => {
    test('', () => {
      const v = mm.i32(2147483647);
      const v2 = mm.i32(-2147483648);
    });
  });

  describe('int64', () => {
    test('', () => {
      const v = mm.i64(9223372036854775807n);
      const v2 = mm.i64(-9223372036854775808n);
    });
  });

  describe('uint', () => {
    test('', () => {
      const v = mm.u(42n);
      const v2 = mm.u(0n);
    });
  });

  describe('uint8', () => {
    test('', () => {
      const v = mm.u8(255);
      const v2 = mm.u8(0);
    });
  });

  describe('uint16', () => {
    test('', () => {
      const v = mm.u16(65535);
      const v2 = mm.u16(0);
    });
  });

  describe('uint32', () => {
    test('', () => {
      const v = mm.u32(4294967295);
      const v2 = mm.u32(0);
    });
  });

  describe('uint64', () => {
    test('', () => {
      const v = mm.u64(18446744073709551615n);
      const v2 = mm.u64(0n);
    });
  });

  describe('float32', () => {
    test('', () => {
      const v = mm.f32(3.14);
      const v2 = mm.f32(0.0);
      const v3 = mm.f32(-1.5);
    });
  });

  describe('float64', () => {
    test('', () => {
      const v = mm.f64(3.141592653589793);
      const v2 = mm.f64(0.0);
      const v3 = mm.f64(-2.71828);
    });
  });

  describe('bigint', () => {
    test('', () => {
      const v = mm.bigint(9223372036854775807n);
      const v2 = mm.bigint(-9223372036854775808n);
    });
  });

  describe('datetime', () => {
    test('', () => {
      const v = mm.datetime(new Date());
    });
  });

  describe('date', () => {
    test('', () => {
      const v = mm.date(new Date());
    });
  });

  describe('time', () => {
    test('', () => {
      const v = mm.time(new Date());
    });
  });

  describe('uuid', () => {
    test('', () => {
      const v = mm.uuid('123e4567-e89b-12d3-a456-426614174000');
    });
  });

  describe('decimal', () => {
    test('', () => {
      const v = mm.decimal('123.45');
      const v2 = mm.decimal('0.00');
    });
  });

  describe('ip', () => {
    test('', () => {
      const v = mm.ip('192.168.1.1');
      const v2 = mm.ip('2001:0db8:85a3:0000:0000:8a2e:0370:7334');
    });
  });

  describe('url', () => {
    test('', () => {
      const v = mm.url('https://example.com');
    });
  });

  describe('email', () => {
    test('', () => {
      const v = mm.email('test@example.com');
    });
  });

  describe('enum', () => {
    test('', () => {
      const v = mm.enum(0);
      const v2 = mm.enum(1);
    });
  });

  describe('image', () => {
    test('', () => {
      const v = mm.image(new Uint8Array([0x89, 0x50, 0x4e, 0x47]));
    });
  });

  describe('video', () => {
    test('', () => {
      const v = mm.video(new Uint8Array([0x00, 0x00, 0x00, 0x18]));
    });
  });

  describe('doc', () => {
    test('', () => {
      const v = mm.doc({
        title: mm.str('Hello'),
        count: mm.i(42n),
      });
    });
  });
});

describe('Encode Value', () => {
  describe('', () => {
    test('', () => {
      const obj = mm.obj({
        kkjjj: mm.i(1212n),
        '2323': mm.map({
          hhhh: mm.slice([mm.str('hello')]),
        }),
      });
    });
  });
});

describe('mm1', () => {
  describe('', () => {
    test('', () => {
      @mm({ desc: '用户' })
      class User {
        @mm({ type: ValueType.Int64, desc: '用户ID', nullable: false })
        id: bigint = 0n;

        @mm({ desc: '昵称' })
        name: string = '';

        @mm({ type: ValueType.Uint8 })
        age: number = 0;
      }

      const u = new User();
      u.id = 666n;
      u.name = 'abc';
      u.age = 20;

      const ast = toMM(u);
      console.log(ast);
    });
  });
});

describe('value to node', () => {
  describe('', () => {
    test('', () => {
      @mm({ desc: '用户' })
      class User {
        @mm({ type: ValueType.Int64, desc: '用户ID', nullable: false })
        id: bigint = 0n;

        @mm({ desc: '昵称' })
        name: string = '';

        @mm({ type: ValueType.Uint8 })
        age: number = 0;
      }

      const u = new User();
      u.id = 666n;
      u.name = 'abc';
      u.age = 20;

      const node = ValueToNode(u);
      console.log(node);
    });
  });
});

describe('encode and decode', () => {
  describe('', () => {
    test('', () => {
      @mm({ desc: '用户' })
      class User {
        @mm({ type: ValueType.Int64, desc: '用户ID', nullable: false })
        id: bigint = 0n;
        @mm({ desc: '昵称' })
        name: string = '';
        @mm({ type: ValueType.Uint8 })
        age: number = 0;
      }

      const u = new User();
      u.id = 666n;
      u.name = 'abc';
      u.age = 20;

      // const node = ValueToNode(u);
      // console.log('ValueToNode', node);
      const wire = encodeFromValue(u);
      console.log('wire', wire);
      const u2 = decodeToValue(wire, User);

      console.log(u2);
    });
  });
});
