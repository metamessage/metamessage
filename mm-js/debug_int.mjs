import { MMEncoder } from './src/mm/encoder.js';
import { MMDecoder } from './src/mm/decoder.js';
import { intLen, MMPrefix, MMConstants } from './src/mm/constants.js';

function debugInt(value) {
  console.log(`\n=== Testing value: ${value}n ===`);

  const encoder = new MMEncoder();
  encoder.encodeUInt(value);
  const bytes = encoder.buffer.data;

  console.log('Encoded bytes:', bytes.map(b => '0x' + b.toString(16).padStart(2, '0')).join(' '));
  console.log('First byte: 0x' + bytes[0].toString(16));

  const firstByte = bytes[0];
  const prefix = firstByte & MMConstants.PrefixMask;
  console.log('Prefix:', prefix, '(PositiveInt expected:', MMPrefix.PositiveInt, ')');

  const suffix = firstByte & MMConstants.IntLenMask;
  console.log('Suffix (len bits):', suffix, '(IntLen2Byte expected:', MMConstants.IntLen2Byte, ')');

  const { extraBytes, len } = intLen(firstByte);
  console.log('intLen result:', { extraBytes, len });

  const decoder = new MMDecoder(bytes);
  const result = decoder.decode();

  console.log('Decoded value:', result.value);
  console.log('Expected:', value);
  console.log('Match:', result.value === value ? 'YES' : 'NO');
}

debugInt(123456n);
debugInt(0n);
debugInt(23n);
debugInt(255n);
debugInt(256n);
debugInt(65535n);
debugInt(65536n);