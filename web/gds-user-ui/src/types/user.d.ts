interface IVasp {
  id: string;
  created_at?: string;
  domain?: string;
  name?: string;
  refresh_token?: string;
}

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
  vasp?: IVasp;
  role?: string,
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

type TUserCollaboratorPermission =
  | 'read:collaborators'
  | 'create:collaborators'
  | 'update:collaborators'
  | 'approve:collaborators';
type TUserCertificatePermission =
  | 'read:certificates'
  | 'create:certificates'
  | 'update:certificates'
  | 'revoke:certificates';
type TVaspPermission = 'read:vasp' | 'create:vasp' | 'update:vasp';
type TOrgnization = 'create:organizations' | 'read:organizations' | 'update:organizations';

type TUserPermission =
  | TUserCollaboratorPermission
  | TUserCertificatePermission
  | TVaspPermission
  | TOrgnization;

type TUserRole = 'Organization Leader' | 'Organization Collaborator';

type TCollaboratorStatus = 'Pending' | 'Confirmed';

type TUserAuthType = 'auth0' | 'google-oauth2' | 'facebook';
