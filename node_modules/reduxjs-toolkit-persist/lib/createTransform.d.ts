declare type TransformConfig = {
    whitelist?: Array<string>;
    blacklist?: Array<string>;
};
export default function createTransform(inbound: Function, outbound: Function, config?: TransformConfig): any;
export {};
