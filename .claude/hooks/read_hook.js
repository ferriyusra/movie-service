async function main() {
    const chunks = [];
    for await (const chunk of process.stdin) {
        chunks.push(chunk);
    }
    const toolArgs = JSON.parse(Buffer.concat(chunks).toString());

    // readPath is the path to the file that Claude is trying to read
    const readPath = toolArgs.tool_input?.file_path || toolArgs.tool_input?.path || "";

    // Allow docs/openapi.yaml but block all other .yaml files
    if (readPath.endsWith(".yaml")) {
        // Normalize the path to handle different path formats
        const normalizedPath = readPath.replace(/\\/g, '/');
        
        // Check if this is the allowed openapi.yaml file
        if (!normalizedPath.endsWith("docs/openapi.yaml")) {
            console.error("You cannot read .yaml files except for docs/openapi.yaml");
            process.exit(2);
        }
    }
}

main();