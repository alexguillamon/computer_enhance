bits 16

add bx, [bx + si]
add bx, [bp + 0]
add si, byte 2
add bp, byte 2
add cx, byte 8
add bx, [bp + 0]
add cx, [bx + 2]
add bh, [bp + si + 4]
add di, [bp + di + 6]
add [bx + si], bx
add [bp + 0], bx
add [bp + 0], bx
add [bx + 2], cx
add [bp + si + 4], bh
add [bp + di + 6], di
add [bx], byte 34
add [bp + si + 1000], byte 29
add ax, [bp + 0]
add al, [bx + si]
add ax, bx
add al, ah
add ax, 1000
add al, -30
add al, 9
sub bx, [bx + si]
sub bx, [bp + 0]
sub si, byte 2
sub bp, byte 2
sub cx, byte 8
sub bx, [bp + 0]
sub cx, [bx + 2]
sub bh, [bp + si + 4]
sub di, [bp + di + 6]
sub [bx + si], bx
sub [bp + 0], bx
sub [bp + 0], bx
sub [bx + 2], cx
sub [bp + si + 4], bh
sub [bp + di + 6], di
sub [bx], byte 34
sub [bx + di], byte 29
sub ax, [bp + 0]
sub al, [bx + si]
sub ax, bx
sub al, ah
sub ax, 1000
sub al, -30
sub al, 9
cmp bx, [bx + si]
cmp bx, [bp + 0]
cmp si, byte 2
cmp bp, byte 2
cmp cx, byte 8
cmp bx, [bp + 0]
cmp cx, [bx + 2]
cmp bh, [bp + si + 4]
cmp di, [bp + di + 6]
cmp [bx + si], bx
cmp [bp + 0], bx
cmp [bp + 0], bx
cmp [bx + 2], cx
cmp [bp + si + 4], bh
cmp [bp + di + 6], di
cmp [bx], byte 34
cmp [4834], byte 29
cmp ax, [bp + 0]
cmp al, [bx + si]
cmp ax, bx
cmp al, ah
cmp ax, 1000
cmp al, -30
cmp al, 9
jne 2
jne -4
jne -6
jne -4
je -2
jl -4
jle -6
jb -8
jbe -10
jp -12
jo -14
js -16
jne -18
jnl -20
jnle -22
jnb -24
jnbe -26
jnp -28
jno -30
jns -32
loop -34
loopz -36
loopnz -38
jcxz -40