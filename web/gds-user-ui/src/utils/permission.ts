import Store from 'application/store';
import { isArray } from 'lodash';
import { USER_PERMISSION } from 'types/enums';

export const getUserPermissionFromStore = () => {
  return Store.getState()?.user?.user?.permissions;
};

/**  hasPermission function
 *   @params permission: string | string[]
 *   @return boolean
 */

export const hasPermission = (permission: TUserPermission | TUserPermission[]) => {
  const userPermission = getUserPermissionFromStore();
  if (isArray(permission)) {
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
    return role.every((r) => userRole.includes(r));
  }

  return userRole?.includes(role);
};


export const canInviteCollaborator = () => {
  return hasPermission(USER_PERMISSION.UPDATE_COLLABORATOR);
};

export const canCreateOrganization = () => {
  return hasPermission(USER_PERMISSION.CREATE_ORGANIZATIONS);
};

