export type IAccount = {
  id_code: string;
  email: string,
  phoneNumber: string,
  jobTitle: string,
  officeAddress: string,
  firstName: string,
  middleName: string,
  lastName: string,
  sex: string,
  gender: string,
  age: string,
  height: string,
  homeAddress: string,
  birthdate: string,
  birthplace: string,
  securityAccessLevel: string
}

export type INewAccount = {
  email: string,
  password: string,
  phone_number: string,
  job_title: string,
  office_address: string,
  employee_identity_data: 
  {
    first_name: string,
    middle_name: string,
    last_name: string,
    sex: string,
    gender: string,
    age: string,
    height: string,
    home_address: string,
    birthdate: string,
    birthplace: string,
  },
  security_access_level: string
}

export type INewDocument = {
  title: string,
  description: string,
  security_access_level: string,
  language: string,
}

export type IUpdateDocument = {
  document_id: string;
  title: string,
  description: string,
  security_access_level: string,
  language: string,
}

export type IContextType = {
  account: IAccount;
  isLoading: boolean;
  setAccount: React.Dispatch<React.SetStateAction<IAccount>>;
  isAuthenticated: boolean,
  setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
  checkAuthAccount: () => Promise<boolean>
}


/************************* */
export type INavLink = {
  imgURL: string;
  route: string;
  label: string;
};

export type IUpdateUser = {
  userId: string;
  name: string;
  bio: string;
  imageId: string;
  imageUrl: URL | string;
  file: File[];
};

export type INewPost = {
  userId: string;
  caption: string;
  file: File[];
  location?: string;
  tags?: string;
};

export type IUpdatePost = {
  postId: string;
  caption: string;
  imageId: string;
  imageUrl: URL;
  file: File[];
  location?: string;
  tags?: string;
};

export type IUser = {
  id: string;
  name: string;
  username: string;
  email: string;
  imageUrl: string;
  bio: string;
};

export type INewUser = {
  name: string;
  email: string;
  username: string;
  password: string;
};
