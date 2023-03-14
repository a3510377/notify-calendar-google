const MaxUint64 = 1n << (64n - 1n);
const all = (min, max, step) => {
  [min, max, step] = [min, max, step].map((d) => BigInt(d));
  let bits = 0n;
  for (let i = min; i <= max; i += step) bits |= 1n << i;
  return bits;
};
