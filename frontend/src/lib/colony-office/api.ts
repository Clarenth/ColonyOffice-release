import { INewAccount, INewDocument, IUpdateDocument } from '@/types'
import { serverConfig } from './config';
import { Query,  } from '@tanstack/react-query';

export async function postCreateAccount(account: INewAccount) {
  const url = serverConfig.signup;
  const payload = 
  {
    method: 'POST',
    headers: 
    {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(
      {
        email: account.email,
        password: account.password,
        phone_number: account.phone_number,
        job_title: account.job_title,
        office_address: account.office_address,
        employee_identity_data: 
        {
          first_name: account.employee_identity_data.first_name,
          middle_name: account.employee_identity_data.middle_name,
          last_name: account.employee_identity_data.last_name,
          sex: account.employee_identity_data.sex,
          gender: account.employee_identity_data.gender,
          age: account.employee_identity_data.age,
          height: account.employee_identity_data.height,
          home_address: account.employee_identity_data.home_address,
          birthdate: account.employee_identity_data.birthdate,
          birthplace: account.employee_identity_data.birthplace,
        },
        security_access_level: account.security_access_level
      }
    )
  }

  try {
    const newAccount = await fetch(url, payload)
    console.log(newAccount)

    if(!newAccount) throw Error;
    return newAccount;
  } catch (error) {
    console.log(error)
    return error;
  }
}

export async function postLoginAccount(account: { email: string, password: string }) {
  const url = serverConfig.login;
  const payload = 
  {
    method: 'POST',
    headers: 
    {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(
      {
        email: account.email,
        password: account.password,
      }
    )
  }

  try {
    const login = await fetch(url, payload)
    .then(response => response.json())
    .then(data => {
      sessionStorage.setItem("idToken", data.tokens.idToken)
      sessionStorage.setItem("refreshToken", data.tokens.refreshToken)
      /*
      Check if what comes back matches a JWT structure
      We should not just allow anything to be placed in the session storage
      */
     return true
    })
    .catch(error => console.log(error))
    //console.log(login)

    if (!login) throw Error;
    return login;
  } catch (error) {
    console.log(error)
    return error;
  }
}

export async function postLogoutAccount() {
  const url = serverConfig.logout;
  try {
    //await fetch(url, 
    const logout = await fetch(url, 
      {
        method: 'POST',
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`
        }
      }
    )
    .then(response => response.json())
    .then(data => {
      sessionStorage.removeItem("idToken");
      sessionStorage.removeItem("refreshToken");
      console.log(data)
    })
    return logout

  } catch (error) {
    console.log(error)
  }
}

export async function getCurrentAccount() {
  const url = serverConfig.currentAccount;
  try {
    const fetchAccount = await fetch(url,
      {
        method: 'GET',
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`
        }
      }
    )
    .then(response => response.json())
    .then((data) => {
      const account = data;
      return account;
    })
    return fetchAccount;
  } catch (error) {
    console.log(error)
  }
}

export async function getAccountPromiseAll() {
  const url = serverConfig.currentAccount;
  const getAccount = await (await fetch(url, {
    method: 'GET',
    headers:
    {
      Authorization: `Bearer ${sessionStorage.getItem("idToken")}`
    }
  })).json()

  const data = Promise.all(
    [getAccount].map((obj) => {
      console.log(obj)
      const account = obj
      /*to get a value from obj: 
        obj.account.employee_identity_data.<field>,
      */
      return account
    })
  )
  return data
}

/********** Tokens Mutations **********/
export async function postNewTokenPair() {
  const url = serverConfig.newTokenPair;
  const refreshToken = sessionStorage.getItem("refreshToken")
  const payload =
  {
    method: 'POST',
    headers: 
    {
      Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
      "Content-Type": "application/json",

    },
    body: JSON.stringify(
      {
        refreshToken: refreshToken,
      }
    )
  }
  try {
    //const tokens = await fetch(url, payload)
    const newTokens = await fetch(url, payload)
    .then(response => response.json())
    .then((data) => {
      console.log(data)
      sessionStorage.setItem("idToken", data.tokens.idToken)
      sessionStorage.setItem("refreshToken", data.tokens.refreshToken)
      /*
      Check if what comes back matches a JWT structure
      We should not just allow anything to be placed in the session storage
      */
     return true
    })
    //if (!newTokens) throw Error;
    return newTokens;
    //return newTokens
  } catch (error) {
    console.log(error)
    return false
  }
}

/********** Documents Mutations **********/
export async function postCreateDocument(document: INewDocument) {
  const url = serverConfig.createDocument;
  const payload = 
  {
    method: 'POST',
    headers:
    {
      Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
      "Content-Type": "application/json",

    },
    body: JSON.stringify(
    {
      title: document.title,
      description: document.description,
      language: document.language,
      security_access_level: document.security_access_level,
    })
  }
  try {
    const postDocument = await fetch(url, payload)
    .then(response => response.json())
    .then((data) => {
      const documentID = data.document_code;
      console.log(data)
      return documentID;
    })
    return postDocument;
  } catch (error) {
    console.log(error)
  }
}

export async function patchUpdateDocument(document: IUpdateDocument) {
  const url = serverConfig.updateDocument;
  const payload = 
  {
    method: 'PATCH',
    headers:
    {
      Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
      "Content-Type": "application/json",

    },
    body: JSON.stringify(
    {
      title: document.title,
      description: document.description,
      language: document.language,
      security_access_level: document.security_access_level,
    })
  }
  try {
    const updateDocument = await fetch(url, payload)
    .then(response => response.json())
    .then((data) => {
      const documentID = data.document_code;
      console.log(data)
      return documentID;
    })
    return updateDocument;
  } catch (error) {
    console.log(error)
  }
}

export async function deleteDocument(document_id: string) {
  const url = `${serverConfig.deleteDocument}${document_id}`
  const payload = 
  {
    method: 'DELETE',
    headers: 
    {
      Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
    }
  }
  try {
    await fetch(url, payload)
      .then(response => response)
  } catch (error) {
    console.log(error)
  }
}

export async function getDocuments() {
  try {
    const url = serverConfig.getDocuments;
    const fetchDocuments = await fetch(url,
      {
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
        }
      }
    )
    .then(response => response.json())
    .then((data) => {
      console.log(data)
      const docsObj = data;
      return docsObj
    })
    return fetchDocuments
  } catch (error) {
    console.log(error)
  }
}

export async function getDocumentByID(document_id: string) {
  try {
    const url = `${serverConfig.getDocumentByID}${document_id}`
    console.log(url)

    const fetchDocument = await fetch(url,
      {
        method: 'GET',
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
        }
      }
    )
    .then(response => response.json())
    .then((data) => {
      console.log(data)
      const docsObj = data;
      return docsObj
    })
    if(!fetchDocument) throw Error;
    return fetchDocument;
  } catch (error) {
    console.log(error)
  }
}

export async function getRecentDocuments() {
  try {
    const url = serverConfig.getDocuments;
    const fetchDocuments = await fetch(url,
      {
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
        }
      }
    )
    .then(response => response.json())
    .then((data) => {
      console.log(data)
      const docsObj = data;
      return docsObj
    })
    if(!fetchDocuments) throw Error;
    return fetchDocuments

  } catch (error) {
    console.log(error)
  }
}

export async function getPaginationDocs(pageIndex, pageCount) {
  const url = `${serverConfig.getMostRecentDocuments}?page=${pageIndex}&count=${pageCount}`;
  // const url2 = `http://localhost:4000/api/v1/docs?page=${pageIndex}&count=${pageCount}`
  try {
    const result = await fetch(url, 
      {
        method: 'GET',
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
        }
      }
    )
    .then(response => response.json())
    return result;
  } catch (error) {
    console.log(error)
  }
}

/********** Files Mutations **********/
export async function postUploadFile(formData: FormData) {
  const url = serverConfig.uploadFiles;
  try {
    const postFiles = await fetch(url,
      {
        method: 'POST',
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
          "Content-Type": "multipart/form-data",
        },
        body: formData
      }
    )
    .then(response => response.json())
    .then((data) => {
      console.log(data)
      return true
    })
    if(!postFiles) throw new Error;
    return postFiles;
  } catch (error) {
    console.log(error)
  }
}

export async function deleteFile(file_id: string) {
  const url = `${serverConfig.deleteFile}${file_id}`
  const payload = 
  {
    method: 'DELETE',
    headers: 
    {
      Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
    }
  }
  try {
    await fetch(url, payload)
      .then(response => response)
  } catch (error) {
    console.log(error)
  }
}

export async function getFilesByID(files_id: string) {
  try {
    const url = `${serverConfig.getFilesByID}${files_id}`
    console.log(url)

    const fetchDocument = await fetch(url,
      {
        method: 'GET',
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
        }
      }
    )
    .then(response => response.json())
    .then((data) => {
      console.log(data)
      const filesObj = data;
      return filesObj
    })
    if(!fetchDocument) throw Error;
    return fetchDocument;
  } catch (error) {
    console.log(error)
  }
}

export async function getFilesByDocumentID(document_id: string) {
  try {
    const url = `${serverConfig.getFilesByDocumentID}${document_id}`
    console.log(url)

    const fetchDocument = await fetch(url,
      {
        method: 'GET',
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`,
        }
      }
    )
    .then(response => response.json())
    .then((data) => {
      console.log(data)
      const filesObj = data;
      return filesObj
    })
    if(!fetchDocument) throw Error;
    return fetchDocument;
  } catch (error) {
    console.log(error)
  }
}