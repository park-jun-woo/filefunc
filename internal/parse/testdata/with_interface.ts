interface Config {
    host: string;
    port: number;
}

type Options = Partial<Config>;

function createConfig(host: string, port: number): Config {
    return { host, port };
}
