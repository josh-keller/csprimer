section .text
global index
index:
	; rdi: matrix
	; esi: rows
	; edx: cols
	; ecx: rindex
	; r8d: cindex
	imul ecx, edx
	add ecx, r8d
	mov rax, [rdi + rcx*4]

	ret
