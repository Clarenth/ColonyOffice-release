export const serverConfig = {
  currentAccount: import.meta.env.VITE_SERVER_CURRENT_ACCOUNT,
  host: import.meta.env.VITE_SERVER,
  
  // Account
  login: import.meta.env.VITE_SERVER_LOGIN,
  logout: import .meta.env.VITE_SERVER_LOGOUT,
  signup: import.meta.env.VITE_SERVER_SIGNUP,
  newTokenPair: import.meta.env.VITE_SERVER_NEW_TOKEN_PAIR,

  // Documents
  createDocument: import.meta.env.VITE_SERVER_CREATE_DOCUMENTS,
  deleteDocument: import.meta.env.VITE_SERVER_DELETE_DOCUMENT,
  updateDocument: import.meta.env.VITE_SERVER_UPDATE_DOCUMENT,
  getDocuments: import.meta.env.VITE_SERVER_GET_DOCUMENTS,
  getDocumentByID: import.meta.env.VITE_SERVER_GET_DOCUMENTS_BY_ID,
  //Pagination urls
  getMostRecentDocuments: import .meta.env.VITE_SERVER_GET_MOST_RECENT_DOCUMENTS,
  
  // Files
  uploadFiles: import.meta.env.VITE_SERVER_UPLOAD_FILES,
  deleteFile: import.meta.env.VITE_SERVER_DELETE_FILE,
  getFilesByID: import.meta.env.VITE_SERVER_GET_FILE_BY_ID,
  getFilesByDocumentID: import.meta.env.VITE_SERVER_GET_FILE_BY_DOC_ID,

}