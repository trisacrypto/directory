
import Store from 'application/store';
import { AUTH0_TYPE } from 'utils/constants';

export const getUserAuthFromStore = () => {
    return Store.getState()?.user?.user?.authType;
};

export const isSocialLogin = () => getUserAuthFromStore() !== AUTH0_TYPE.AUTH0;
