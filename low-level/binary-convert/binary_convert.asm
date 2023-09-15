section .text
global binary_convert
binary_convert:
	xor		eax, eax

_loop:
	mov	 		cl, [rdi]  ; load current character
	cmp			cl, 0			 ; break loop if end of string
	je			_end
	add			rdi, 1		 ; increment address

	shl			eax, 1		 ; shift current sum left
	sub			cl,	48     ; convert ascii "1" or "0" to 1 or 0
	add			al, cl		 ; add it to the total
	jmp			_loop

_end:
	ret


