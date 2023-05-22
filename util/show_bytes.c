#include <stdio.h>

typedef unsigned char *byte_pointer;

void show_bytes(byte_pointer start, size_t len) {
  int i;
  for (i = 0; i < len; i++) {
    printf(" %.2x", start[i]);
  }
  printf("\n");
}

void print_byte_binary(char byte) {
  for (int i = 0x80; i > 0; i >>= 1) {
    putchar((byte & i) ? '1' : '0');
  }
}

void show_binary(byte_pointer start, size_t len) {
  int i;

  for (i = 0; i < len; i++) {
    print_byte_binary(start[i]);
    if (i != len - 1) {
      printf(" ");
    }
  }
  printf("\n");
}


void show_int(int x) {
  show_bytes((byte_pointer) &x, sizeof(int));
}

void show_float(float x) {
  show_bytes((byte_pointer) &x, sizeof(float));
}

void show_pointer(void *x) {
  show_bytes((byte_pointer) &x, sizeof(void *));
}

void show_int_binary(int x) {
  show_binary((byte_pointer) &x, sizeof(int));
}

void show_float_binary(float x) {
  show_binary((byte_pointer) &x, sizeof(float));
}

void show_pointer_binary(void *x) {
  show_binary((byte_pointer) &x, sizeof(void *));
}

void test_show_bytes(int val) {
  int ival = val;
  float fval = (float) ival;
  int *pval = &ival;
  show_int(ival);
  show_float(fval);
  show_pointer(pval);
}

void test_show_binary(int val) {
  int ival = val;
  float fval = (float) ival;
  int *pval = &ival;
  show_int_binary(ival);
  show_float_binary(fval);
  show_pointer_binary(pval);
}

int main() {
  test_show_bytes(12345);
  printf("----------\n");
  test_show_binary(12345);
}
