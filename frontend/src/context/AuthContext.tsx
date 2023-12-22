// Libraries
import React, { createContext, useContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

// API
import { getCurrentAccount, postNewTokenPair } from '@/lib/colony-office/api';

// Types
import { IAccount } from '@/types';

export const INITIAL_ACCOUNT = {
  id_code: "",
  email: "",
  phoneNumber: "",
  jobTitle: "",
  officeAddress: "",
  firstName: "",
  middleName: "",
  lastName: "",
  sex: "",
  gender: "",
  age: "",
  height: "",
  homeAddress: "",
  birthdate: "",
  birthplace: "",
  securityAccessLevel: ""
};

const INITIAL_STATE = {
  account: INITIAL_ACCOUNT,
  jwt: "",
  isLoading: false,
  isAuthenticated: false,
  setAccount: () => {},
  setIsAuthenticated: () => {},
  checkAuthAccount: async () => false as boolean,
  fetchNewTokenPair: async () => false as boolean,
};

export type IContextType = {
  account: IAccount;
  jwt: string;
  isLoading: boolean;
  setAccount: React.Dispatch<React.SetStateAction<IAccount>>;
  isAuthenticated: boolean,
  setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
  checkAuthAccount: () => Promise<boolean>
  fetchNewTokenPair: () => Promise<boolean | undefined>
}

const AuthContext = createContext<IContextType>(INITIAL_STATE);

export function AuthProvider ({ children }: { children: React.ReactNode }) {
  const [account, setAccount] = useState<IAccount>(INITIAL_ACCOUNT)
  const [jwt, setJWT] = useState("")
  const [isLoading, setIsLoading] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const navigate = useNavigate();
  
  const checkAuthAccount = async () => {
    try {
      const currentAccount = await getCurrentAccount();

      if(currentAccount) {
        setAccount({
          id_code: currentAccount.account.id_code,
          email: currentAccount.account.email,
          phoneNumber: currentAccount.account.phoneNumber,
          jobTitle: currentAccount.account.jobTitle,
          officeAddress: currentAccount.account.officeAddress,
          firstName: currentAccount.account.employee_identity_data.first_name,
          middleName: currentAccount.account.employee_identity_data.middle_name,
          lastName: currentAccount.account.employee_identity_data.last_name,
          sex: currentAccount.account.employee_identity_data.sex,
          gender: currentAccount.account.employee_identity_data.gender,
          age: currentAccount.account.employee_identity_data.age,
          height: currentAccount.account.employee_identity_data.height,
          homeAddress: currentAccount.account.employee_identity_data.home_address,
          birthdate: currentAccount.account.employee_identity_data.birthplace,
          birthplace: currentAccount.account.employee_identity_data.burthplace,
          securityAccessLevel: currentAccount.account.employee_identity_data.security_access_level,
        })
        setIsAuthenticated(true);
        return true;
      }
      return false;
    } catch (error) {
      console.log(error)
      return false
    } finally {
      setIsLoading(false);
    }
  };

  const fetchNewTokenPair = async () => {
    try {
      const newTokenPair = await postNewTokenPair();
      console.log("Fired useEffect")

      if(!newTokenPair){
        setIsAuthenticated(true)
        return true
      }
      return false
    } catch (error) {
      console.log(error)
    }
  }

  useEffect(() => {
    fetchNewTokenPair();
    if
    (
      sessionStorage.getItem("idToken") === "[]" ||
      //sessionStorage.getItem("idToken") === null ||
      sessionStorage.getItem("refreshToken") === "[]"
      //sessionStorage.getItem("refreshToken") === null
    ) navigate("/login")
    checkAuthAccount();
  }, [])

  const value = {
    account,
    jwt,
    setAccount,
    isLoading,
    isAuthenticated,
    setIsAuthenticated,
    checkAuthAccount,
    fetchNewTokenPair,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

//export default AuthProvider;

export const useAccountContext = () => useContext(AuthContext)