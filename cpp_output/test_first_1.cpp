#include <bits/stdc++.h>

template <typename T>
	void _format_output(std::ostream& out, const T& str) 
	{	
		out << str;
	}

auto main() -> int {
    // Declare variables
int x;
int y;

    // Assign values
x = 10;
y = 20;

    // Perform arithmetic operations
auto sum = x + y;
auto difference = x - y;
auto product = x * y;
auto quotient = x / y;

    // Print results
std::cout << "Sum:" << " ";
_format_output(std::cout, sum);
std::cout << std::endl;
std::cout << "Difference:" << " ";
_format_output(std::cout, difference);
std::cout << std::endl;
std::cout << "Product:" << " ";
_format_output(std::cout, product);
std::cout << std::endl;
std::cout << "Quotient:" << " ";
_format_output(std::cout, quotient);
std::cout << std::endl;
return 0;
}

