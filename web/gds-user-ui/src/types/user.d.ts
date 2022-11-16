type TUserPermission = 'collaborators' | 'certificate' | 'admin';
interface IUserState {
  name: string;
  email: string;
  roles: string[];
  pictureUrl: string;
  permission?: TUserPermission;
}
type TUser = {
  isFetching?: boolean;
  isSuccess?: boolean;
  isError?: boolean;
  errorMessage?: string;
  isLoggedIn: boolean;
  user: IUserState | null;

};
