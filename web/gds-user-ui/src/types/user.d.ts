
interface IUserState {
  name: string;
  email: string;
  roles: string[];
  pictureUrl: string;
  permissions?: TUserPermission[];
}
type TUser = {
  isFetching?: boolean;
  isSuccess?: boolean;
  isError?: boolean;
  errorMessage?: string;
  isLoggedIn: boolean;
  user: IUserState | null;

};

type TUserCollaboratorPermission = 'read:collaborators' | 'create:collaborators' | 'update:collaborators' | 'approve:collaborators';
type TUserCertificatePermission = 'read:certificates' | 'create:certificates' | 'update:certificates' | 'revoke:certificates';
type TVaspPermission = 'read:vasp' | 'create:vasp' | 'update:vasp';

type TUserPermission = TUserCollaboratorPermission | TUserCertificatePermission | TVaspPermission;

type TUserRole = 'Organization Leader' | 'Organization Collaborator';

type TCollaboratorStatus = 'pending' | 'joined';

