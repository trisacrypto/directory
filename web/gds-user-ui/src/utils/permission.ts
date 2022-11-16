

import Store from 'application/store';
import { isArray } from 'lodash';

export const getUserPermissionFromStore = () => {
    return Store.getState()?.user?.user?.permission;
};

/**  hasPermission function
 *   @params permission: string | string[]
 *   @return boolean
*/

export const hasPermission = (permission: TUserPermission | TUserPermission[]) => {
    const userPermission = getUserPermissionFromStore();
    if (isArray(permission)) {
        // all permission element should be in userPermission
        return permission.every((p) => userPermission.includes(p));
    }

    return userPermission?.includes(permission);
};

/**  hasRole function
 *   @params role: string | string[]
 *   @return boolean
*/

export const hasRole = (role: TUserRole | TUserRole[]) => {
    const userRole = Store.getState()?.user?.user?.role;
    if (isArray(role)) {
        // all role element should be in userRole
        return role.every((r) => userRole.includes(r));
    }

    return userRole?.includes(role);
};
