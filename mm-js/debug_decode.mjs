import { MMEncoder } from './src/mm/encoder.js';
import { MMDecoder } from './src/mm/decoder.js';
import { intLen, MMPrefix, MMConstants } from './src/mm/constants.js';

function debugDecode(value) {
  console.log(`\n=== Decoding ${value}n ===`);

  const encoder = new MMEncoder();
  encoder.encodeUInt(value);
  const bytes = encoder.buffer.data;

  console.log('Encoded bytes:', bytes.map(b => '0x' + b.toString(16).padStart(2, '0')).join(' '));
  console.log('Bytes buffer length:', bytes.length);

  const decoder = new MMDecoder(bytes);
  console.log('Decoder buffer:', decoder.buffer.buffer.map(b => '0x' + b.toString(16).padStart(2, '0')).join(' '));
  console.log('Decoder offset:', decoder.buffer.offset);

  const firstByte = decoder.buffer.read();
  console.log('First byte read:', '0x' + firstByte.toString(16));

  const prefix = firstByte & MMConstants.PrefixMask;
  console.log('Prefix:', prefix, '(PositiveInt:', MMPrefix.PositiveInt, ')');

  const suffix = firstByte & MMConstants.IntLenMask;
  console.log('Suffix (len bits):', suffix);

  const lenResult = intLen(firstByte);
  console.log('intLen result:', lenResult);

  // Re-create decoder to test full decode
  const decoder2 = new MMDecoder(bytes);
  const result = decoder2.decode();
  console.log('Decoded:', result.value, '(expected:', value, ')');
  console.log('Match:', result.value === value);
}

debugDecode(123456n);
debugDecode(0n);
debugDecode(23n);
debugDecode(255n);
debugDecode(256n);
debugDecode(65535n);
debugDecode(65536n);