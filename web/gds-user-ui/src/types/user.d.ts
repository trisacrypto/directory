
interface IUserState {
  name: string;
  email: string;
  roles: string[];
  pictureUrl: string;
  permissions?: TUserPermission[];
  authType?: string; // auth0, google-oauth,facebook, etc for now we use sub key of idTokenPayload to identify auth type
  id?: string;
  lastLogin?: string;
  createdAt?: string;
  hasSetPassword?: boolean; // this is for user who sign up with social account who need to set password
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

type TUserAuthType = 'auth0' | 'google-oauth2' | 'facebook';

