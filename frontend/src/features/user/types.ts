export type UpdateProfileData = {
  name?: string;
  avatar_url?: string;
};

export type ChangePasswordData = {
  old_password: string;
  new_password: string;
};
