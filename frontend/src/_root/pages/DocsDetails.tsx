// Libraries
import React from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

// Mutations
import { useGetDocumentByIDMutation, useGetFilesByDocumentIDMutation } from '@/lib/react-query/queriesAndMutations'

// Components
import Loader from '@/components/shared/Loader'
import { useAccountContext } from '@/context/AuthContext'
import { Button } from '@/components/ui/button'

const DocsDetails = () => {
  const { account } = useAccountContext();
  const { document_id } = useParams()
  const { data: document, isLoadingDoc } = useGetDocumentByIDMutation(document_id || '')
  const { data: files, isLoadingFiles } = useGetFilesByDocumentIDMutation(document_id || '')
  const navigate = useNavigate();

  const handleDeleteDocument = () => {
    
  }

  return (
    <div className='document_details-container'>
      <div className='hidden md:flex max-w-5x1 w-full'>
        <Button
          onClick={() => navigate(-1)}
          variant="ghost"
          className='shad-button_ghost'>
            <img 
              src='/assets/arrow-left-dark-svgrepo-com.svg'
              width={35}
              height={35}
            />
        </Button>
      </div>
      { isLoadingDoc || !document ? (
        <Loader />
      ) : 
      (
      <div className='document_details-card'>
        <div className='flex items-center w-full'>

          <div className='flex flex-col w-full'>
            <div>
              <div className='relative flex pt-3 items-center'>
                <div className='flex-grow border-t border-gray-400'></div>
                <span className='flex-shrink mx-4 base-medium lg:body-bold'>Title</span>
                <div className='flex-grow border-t border-gray-400'></div>
              </div>
              <div className='pt-1 pl-4'>
                <span>{document?.document?.document_title}</span>
              </div>
            </div>
            

            <div>
              <div className='relative flex pt-3 items-center'>
                <div className='flex-grow border-t border-gray-400'></div>
                <span className='flex-shrink mx-4 base-medium lg:body-bold'>Author</span>
                <div className='flex-grow border-t border-gray-400'></div>
              </div>
              <div className='pt-1 pl-4'>
                <span>{document?.document?.author_name}</span>
              </div>
            </div>

            <div>
              <div className='relative flex pt-3 items-center'>
                <div className='flex-grow border-t border-gray-400'></div>
                <span className='flex-shrink mx-4 base-medium lg:body-bold'>Date</span>
                <div className='flex-grow border-t border-gray-400'></div>
              </div>
              <div className='pt-1 pl-4'>
                <span>{document?.document?.created_at}</span>
              </div>
            </div>

            <div>
              <div className='relative flex pt-3 items-center'>
                <div className='flex-grow border-t border-gray-400'></div>
                <span className='flex-shrink mx-4 base-medium lg:body-bold'>Language</span>
                <div className='flex-grow border-t border-gray-400'></div>
              </div>
              <div className='pt-1 pl-4'>
                  <span>{document?.document?.language}</span>
              </div>
            </div>

            <div>
              <div className='relative flex pt-3 items-center'>
                <div className='flex-grow border-t border-gray-400'></div>
                <span className='flex-shrink mx-4 base-medium lg:body-bold'>Description</span>
                <div className='flex-grow border-t border-gray-400'></div>
              </div>
              <div className='pt-1 pl-4'>
                <span>{document?.document?.description}</span>
              </div>
            </div>

            <div>
              <div className='relative flex pt-3 items-center'>
                <div className='flex-grow border-t border-gray-400'></div>
                <span className='flex-shrink mx-4 base-medium lg:body-bold'>Files</span>
                <div className='flex-grow border-t border-gray-400'></div>
              </div>
              <div className='pt-1 px-4'>
                  <span>{ files?.file_title }</span>
                  { isLoadingFiles && !files ? (
                    <Loader />
                    ) : (
                      <ul>
                        { files?.files.map((file) => (
                          <li key={file?.file_id} className='grid p-2 w-full'>
                            <div className='grid grid-cols-2 '>
                              <div className='grid justify-start'>
                                <p>File: {file.file_title}</p>
                                <p>Author(s): {file.author_name}</p>
                                <p>Security Level: {file.security_access_level}</p>
                              </div>
                              <div className='grid justify-end'>
                                <div>
                                <p>Created at: {file.created_at}</p>
                                <p>Updated at: {file.updated_at}</p>
                                </div>
                              </div>
                            </div>
                          </li>
                        ))}
                      </ul>
                    )
                  }
              </div>
            </div>

            <div>
              <div className='relative flex pt-3 items-center'>
                <div className='flex-grow border-t border-gray-400'></div>
                <span className='flex-shrink mx-4 base-medium lg:body-bold'>Actions</span>
                <div className='flex-grow border-t border-gray-400'></div>
              </div>
              <div className='flex-end gap-2'>
                <Button variant="ghost" className='document_details-edit_btn '>
                <Link to={`/edit-docs/${document_id}`} className={`${account.id_code !== document?.document?.author_id && 'hidden'}`}>
                  Edit
                </Link>
                  
                </Button>
                <Button 
                  variant="ghost"
                  className={`document_details-delete_btn ${account.id_code !== document?.document?.author_id && 'hidden'}`}
                  onClick={handleDeleteDocument}>
                    Delete
                </Button>
              </div>
            </div>

          </div>
        </div>
      </div>
      )
      }
    </div>
  )
}

export default DocsDetails