"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function createTransform(
// @NOTE inbound: transform state coming from redux on its way to being serialized and stored
// eslint-disable-next-line @typescript-eslint/ban-types
inbound, 
// @NOTE outbound: transform state coming from storage, on its way to be rehydrated into redux
// eslint-disable-next-line @typescript-eslint/ban-types
outbound, config = {}) {
    const whitelist = config.whitelist || null;
    const blacklist = config.blacklist || null;
    function whitelistBlacklistCheck(key) {
        if (whitelist && whitelist.indexOf(key) === -1)
            return true;
        if (blacklist && blacklist.indexOf(key) !== -1)
            return true;
        return false;
    }
    return {
        in: (state, key, fullState) => !whitelistBlacklistCheck(key) && inbound
            ? inbound(state, key, fullState)
            : state,
        out: (state, key, fullState) => !whitelistBlacklistCheck(key) && outbound
            ? outbound(state, key, fullState)
            : state,
    };
}
exports.default = createTransform;
