import type { PersistConfig, Storage, Transform } from '../types';
declare type V4Config = {
    storage?: Storage;
    serialize: boolean;
    keyPrefix?: string;
    transforms?: Array<Transform<any, any>>;
    blacklist?: Array<string>;
    whitelist?: Array<string>;
};
export default function getStoredState(v4Config: V4Config): (v5Config: PersistConfig<any>) => any;
export {};
