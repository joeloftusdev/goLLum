# goLLum: Copilot Script Generator 
![gollumapp](https://github.com/joeloftusdev/gollum/assets/152509645/a62b2856-c60d-45f9-95c1-69748d66df15)

Gollum is a command-line tool designed to generate scripts in various programming languages. It utilizes the OpenAI GPT-3.5 model to understand prompts and generate corresponding scripts. With Gollum, you can quickly prototype scripts for a variety of purposes without the need to manually write code.

## Features

- Generates scripts based on user prompts
- Supports: Python, Go, Ruby, Perl, Bash, PowerShell, JavaScript, TypeScript, PHP, and Lua.
- Customizable output directory for generated scripts.


## Usage
To use Gollum, you need to have Go installed on your system. Once you have Go installed, you can install Gollum using the following commands:

1. Clone the project
2. Create a text file in the project directory named `apikey.txt` and paste your [OpenAI API key](https://platform.openai.com/api-keys) into it.
3. Install Gollum by executing the following command in your terminal:
```bash
go install
```
4. Run Gollum from anywhere on your pc by executing the following command:
```bash
gollum
```
5. Enter your prompt when prompted, and Gollum will generate the corresponding script for you.

## Example

```bash
Welcome to goLLum. Generate a script or type 'quit' to exit.
You: Write a Python script to calculate factorial of a number. # Keyword Python will generate a .py file
Script generated. Please specify the output directory:
</path/to/your/scripts>
Python script generated in </path/to/your/scripts>
```

## Dependencies

- [Go](golang.org).
- [OpenAI GPT-3.5](https://openai.com): AI model used for goLLum.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

