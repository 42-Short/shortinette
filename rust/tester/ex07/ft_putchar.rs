#![no_std]
#![no_main]
#![no_implicit_prelude]

use core::arch::asm;
use core::panic::PanicInfo;

#[no_mangle]
pub extern "C" fn _start() -> ! {
    ft_putchar(b'4');
    ft_putchar(b'2');
    ft_putchar(b'\n');
    ft_exit(42);
}

fn ft_putchar(c: u8) {
    unsafe {
        asm!(
            // `rax` defines which syscall is to be used (1 = sys_write)
            "mov rax, 1",
            // `rdi` defines the standard output's fd (1 = stdout)
            "mov rdi, 1",
            // `rsi` is a pointer to the buffer to write
            "mov rsi, rdx",
            // `rdx` defines the number of bytes to write
            "mov rdx, 1",
            // triggers the syscall with above params
            "syscall",
            in("rdx") &c,
            lateout("rax") _, lateout("rdi") _, lateout("rsi") _, lateout("rdx") _,
            options(nostack)
        );
    }
}

fn ft_exit(_code: u8) -> ! {
    unsafe {
        asm!(
            "mov rax, 60",
            "syscall",
            options(noreturn)
        );
    }
}

#[panic_handler]
fn panic(_info: &PanicInfo) -> ! {
    loop {}
}
