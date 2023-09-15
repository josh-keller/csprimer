%define MASK 0x07fffffe

section .text
global pangram
pangram:
	; rdi: source string
	; need registers for:
	;  - bitstring (bs)
	;  - offset (can we just use rdi?)
	;  - mask

	; high-level steps:
	;  - first just iterate through the string without segfault
	;  - add ignoring of the first 64 chars
	;  - create the 1 shifted to the correct place (view in LLDB)
	;  - |= with the bs register (check reg in LLDB)
	;  - check mask

	xor		r9d, r9d		; bs = 0

_loop:
  mov		cl, [rdi]     	  ; load current char
	inc		rdi 							; increment offset
	cmp		cl, 0							; if it is zero
	je		_end							; jump to end
	cmp		cl, 0x40					; skip < @
	jl		_loop

	and		cl, 0x1f					; load character mask into cl (look at low order 5 bits)
	lea 	r10d, 1						; load 1 into r10d
	shl		r10d, cl					; shift by cl
	or		r9d, r10d 				; or this with bitstring

	jmp		_loop							; back to beginning of loop


_end:
	and		r9d, MASK		      ; mask the bitstring
	cmp		r9d, MASK		      ; compare the bitstring
	je		_success					; if they are equal go to success
	lea   eax, 0
	ret

_success:
	lea 	eax, 1
	ret
