type TUserPermission = 'collaborators' | 'certificate' | 'admin';
interface IUserState {
  name: string;
  email: string;
  pictureUrl: string;
  permission?: TUserPermission;
}
type TUser = {
  isLoggedIn: boolean;
  user: IUserState | null;
};
