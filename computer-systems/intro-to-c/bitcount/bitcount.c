#include <assert.h>
#include <popcntintrin.h>
#include <stdio.h>
#include <nmmintrin.h>

int bitcount(int);

int main() {
    assert(bitcount(0) == 0);
    assert(bitcount(1) == 1);
    assert(bitcount(3) == 2);
    assert(bitcount(8) == 1);
    // harder case:
    assert(bitcount(0xffffffff) == 32);
    printf("OK\n");
}

int bitcount(int n) {
  return __builtin_popcount(n);
}

// int bitcount(int num) {
//   num = ((num >> 1) & 0x55555555) + (num & 0x55555555);
//   num = ((num >> 2) & 0x33333333) + (num & 0x33333333);
//   num = ((num >> 4) & 0x0f0f0f0f) + (num & 0x0f0f0f0f);
//   num = ((num >> 8) & 0x00ff00ff) + (num & 0x00ff00ff);
//   num = ((num >> 16) & 0x0000ffff) + (num & 0x0000ffff);
//   return num;
// }

// int bitcount(int n) {
//   int count = 0;
//   while (n) {
//     n &= n - 1;
//     count++;
//   }
//
//   return count;
// }

