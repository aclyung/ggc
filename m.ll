%int8 = type i8
%int = type i64
%bool = type i8
%string = type { %int, %int8* }
%float = type float
%main.type.t = type i64
%main.type.matrix = type i64

@true = global %bool 1
@false = global %bool 0
@.str.0 = global [2 x %int8] c"\0A\00"
@.str.1 = global [1 x %int8] c"\00"
@.str.2 = global [1 x %int8] c"\00"

define void @println(...) {
0:
	%1 = getelementptr [1 x %int8], [1 x %int8]* @.str.1, %int 0, %int 0
	%2 = call i32 (%int8*, ...) @printf(%int8* %1)
	%3 = getelementptr [2 x %int8], [2 x %int8]* @.str.0, %int 0, %int 0
	%4 = call i32 (%int8*, ...) @printf(%int8* %3)
	ret void
}

define void @print(...) {
0:
	%1 = getelementptr [1 x %int8], [1 x %int8]* @.str.2, %int 0, %int 0
	%2 = call i32 (%int8*, ...) @printf(%int8* %1)
	ret void
}

declare i32 @printf(%int8* %0, ...)

define void @main.func.hi() {
0:
	%1 = sub %int 1, 1
	ret void
}

define void @main() {
0:
	%1 = alloca [12 x %int8]
	store [12 x %int8] c"3 + 2 = %d\0A\00", [12 x %int8]* %1
	%2 = getelementptr [12 x %int8], [12 x %int8]* %1, %int 0, %int 0
	%3 = add %int 3, 2
	%4 = call i32 (%int8*, ...) @printf(%int8* %2, %int %3)
	%5 = alloca [12 x %int8]
	store [12 x %int8] c"3 - 2 = %d\0A\00", [12 x %int8]* %5
	%6 = getelementptr [12 x %int8], [12 x %int8]* %5, %int 0, %int 0
	%7 = sub %int 3, 2
	%8 = call i32 (%int8*, ...) @printf(%int8* %6, %int %7)
	%9 = alloca [12 x %int8]
	store [12 x %int8] c"1 / 3 = %d\0A\00", [12 x %int8]* %9
	%10 = getelementptr [12 x %int8], [12 x %int8]* %9, %int 0, %int 0
	%11 = sdiv %int 1, 3
	%12 = call i32 (%int8*, ...) @printf(%int8* %10, %int %11)
	%13 = alloca [12 x %int8]
	store [12 x %int8] c"3 * 2 = %d\0A\00", [12 x %int8]* %13
	%14 = getelementptr [12 x %int8], [12 x %int8]* %13, %int 0, %int 0
	%15 = mul %int 3, 2
	%16 = call i32 (%int8*, ...) @printf(%int8* %14, %int %15)
	%17 = alloca [13 x %int8]
	store [13 x %int8] c"3 %% 2 = %d\0A\00", [13 x %int8]* %17
	%18 = getelementptr [13 x %int8], [13 x %int8]* %17, %int 0, %int 0
	%19 = srem %int 3, 2
	%20 = call i32 (%int8*, ...) @printf(%int8* %18, %int %19)
	ret void
}

define %int @main.func.sum() {
0:
	%1 = add %int 1, 1
	ret %int 456
}