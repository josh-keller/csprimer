

						global sum_to_n

						section .text

sum_to_n:   xor 		rax, rax    ; Set sum to 0

loop:				cmp  		rdi, 0			; Check if rdi is zero
						je			return      ; go to return if it is

						add 		rax, rdi    ; add n to sum
						dec     rdi         ; decrement n
						jmp 		loop        ; back to the top of the loop

return:   	ret
; n will be in rdi/edi
; 
; Two registers:
;  - one for sum
;  - one for incrementing number

; Set incrementing number to 0
; 
; Compare incr and n
; If incr > n jump to end (to return)
; Add increment to sum
; Add 1 to increment
; jump to the compare
; 
; return sum
; 
; sum will go in rax/eax
