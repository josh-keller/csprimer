default rel

				section .text
				global 	volume
volume:
        vmulss  xmm0, xmm0, xmm0     ; r^2
				vmulss  xmm0, xmm0, [.PI]    ; * pi
				vdivss  xmm1, xmm1, [.Three] ; h/3
  			vmulss	xmm0,	xmm0, xmm1
 				ret

				section .data
.PI:    dd   3.14159
.Three: dd   3.0
