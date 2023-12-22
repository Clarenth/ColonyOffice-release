import {
  useQuery,
  useQueryClient,
  useMutation,
  useInfiniteQuery,
} from '@tanstack/react-query';

// API
import { 
  postCreateAccount, 
  postLoginAccount, 
  postLogoutAccount,
  postCreateDocument,
  postUploadFile,
  getRecentDocuments,
  getCurrentAccount,
  getDocumentByID,
  getFilesByDocumentID,
  deleteDocument,
  patchUpdateDocument,
} from '../colony-office/api';

// Types
import { INewAccount, INewDocument, IUpdateDocument } from '@/types';
import { QUERY_KEYS } from './querykeys';

/*
Mutations are used as a middleware between the Client functions and the Server-side functions
They can, or are, defined as having two parts: Query and Mutation
-Query: is the Read in CRUD. Use for functions like Get calls, and an SQL Select statement
-Mutation: is the Create, Update, and Delete in CRUD. Use for functions like POST, PATCH, PUT, DELETE
*/

/********** Account Mutations **********/
export const useCreateAccountMutation = () => {
  return useMutation({
    mutationFn: (account: INewAccount) => postCreateAccount(account)
  })
}

export const useLoginAccountMutation = () => {
  return useMutation({
    mutationFn: (account: {
      email: string, 
      password: string
    }) => postLoginAccount(account)
  })
}
 
export const useLogoutAccountMutation = () => {
  return useMutation({
    mutationFn: () => postLogoutAccount()
  })
}

export const useGetCurrentAccount = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.GET_CURRENT_ACCOUNT],
    queryFn: getCurrentAccount,
  })
}

/********** Document Mutations **********/
export const useCreateDocumentMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (document: INewDocument) => postCreateDocument(document),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.GET_RECENT_DOCUMENTS]
      })
    }
  })
}

export const useDeleteDocumentMutation = () => {
return useMutation({
  mutationFn: (document_id: string) => deleteDocument(document_id)
})
}

export const useGetRecentDocumentsMutation = () => {
  return useQuery({
    queryKey: [QUERY_KEYS.GET_RECENT_DOCUMENTS],
    queryFn: getRecentDocuments,
  })
}

export const useGetDocumentByIDMutation = (document_id: string) => {
  return useQuery({
    queryKey: [QUERY_KEYS.GET_DOCUMENT_BY_ID, document_id],
    queryFn: () => getDocumentByID(document_id),
    enabled: !!document_id
  })
}

export const useUpdateDocumentMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (document: IUpdateDocument) => patchUpdateDocument(document),
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.GET_DOCUMENT_BY_ID, data?.document_id]
      })
    }
  })
}

/********** Files Mutations **********/
export const useUploadFileMutation = () => {
  return useMutation({
    mutationFn: (file: FormData) => postUploadFile(file)
  })
}

export const useGetFilesByDocumentIDMutation = (document_id: string) => {
  return useQuery({
    queryKey: [QUERY_KEYS.GET_FILES_BY_DOC_ID, document_id],
    queryFn: () => getFilesByDocumentID(document_id),
    enabled: !!document_id
  })
}