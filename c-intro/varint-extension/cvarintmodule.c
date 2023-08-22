#include <stdint.h>
#include <stdio.h>
#define PY_SSIZE_T_CLEAN
#include <python3.10/Python.h>
#include <python3.10/listobject.h>
#include <python3.10/bytesobject.h>
#include <python3.10/longobject.h>

static PyObject *cvarint_encode(PyObject *self, PyObject *args) {
  int i = 0;
  uint64_t num;
  char result[11];

  if(!PyArg_ParseTuple(args, "k", &num))
    return NULL;

  while (num > 0) {
    result[i] = num & 0x7f;
    num >>= 7;
    result[i] |= (num ? 0x80 : 0x00);
    i++;
  }
  result[i] = 0x00;
  return PyBytes_FromString(result);
}

static PyObject *cvarint_decode(PyObject *self, PyObject *args) {
  uint64_t n = 0;
  int shift = 0;
  Py_ssize_t len;
  const unsigned char *b;
  if(!PyArg_ParseTuple(args, "y#", &b, &len))
    return NULL;

  for (Py_ssize_t i = len - 1; i >= 0; i--) {
    n <<= 7;
    n |= (b[i] & 0x7f);
  }

  return PyLong_FromUnsignedLongLong(n);
}

static PyMethodDef CVarintMethods[] = {
    {"encode", cvarint_encode, METH_VARARGS, "Encode an integer as varint."},
    {"decode", cvarint_decode, METH_VARARGS,
     "Decode varint bytes to an integer."},
    {NULL, NULL, 0, NULL}};

static struct PyModuleDef cvarintmodule = {
    PyModuleDef_HEAD_INIT, "cvarint",
    "A C implementation of protobuf varint encoding", -1, CVarintMethods};

PyMODINIT_FUNC PyInit_cvarint(void) { return PyModule_Create(&cvarintmodule); }
