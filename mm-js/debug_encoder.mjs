import { MMEncoder } from './src/mm/encoder.js';
import { intLen, MMPrefix, MMConstants } from './src/mm/constants.js';

function debugEncoder(value) {
  console.log(`\n=== Encoder Debug: ${value}n ===`);

  const encoder = new MMEncoder();
  encoder.encodeUInt(value);
  const bytes = encoder.buffer.data;

  console.log('Encoded bytes:', bytes.map(b => '0x' + b.toString(16).padStart(2, '0')).join(' '));

  const { extraBytes } = encoder.calcIntExtraBytes(value);
  console.log('calcIntExtraBytes:', { extraBytes });

  const firstByte = bytes[0];
  console.log('First byte: 0x' + firstByte.toString(16));
  console.log('IntLenMask bits:', firstByte & MMConstants.IntLenMask);
  console.log('IntLen1Byte:', MMConstants.IntLen1Byte);
  console.log('IntLen2Byte:', MMConstants.IntLen2Byte);

  const result = intLen(firstByte);
  console.log('intLen result:', result);
}

debugEncoder(0n);
debugEncoder(23n);
debugEncoder(255n);
debugEncoder(256n);
debugEncoder(65535n);
debugEncoder(65536n);
debugEncoder(123456n);