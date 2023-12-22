// Libraries
import React from 'react'
import { Link } from "react-router-dom"

// Components
import DocsNavbar2 from '@/components/shared/DocsNavbar2'
import DocsNavbar from '@/components/shared/DocsNavbar'
import DocsCard from '@/components/shared/DocsCard';
import Loader from '@/components/shared/Loader';

// API
import { useGetRecentDocumentsMutation } from '@/lib/react-query/queriesAndMutations';


const Docs = () => {
  const { data: documents, isPending: isLoadingDocs, isError: isErrorDocs } = useGetRecentDocumentsMutation();


  return (
    <React.Fragment>
      <nav>
        <DocsNavbar />
      </nav>
      <div className='flex flex-1'>
        <div className='home-container'>
          <div className='home-posts'>
            { isLoadingDocs && !documents ? (
              <Loader />
              ) : (
                <ul className='flex flex-col flex-1 gap-9 w-full'>
                  { documents?.documents.map((document) => (
                    <li key={document?.document_id} className='flex justify-center w-full'>
                      <DocsCard document={document}/>
                    </li>
                  ))}
                </ul>
              )
            }
          </div>
        </div>
      </div>
    </React.Fragment>
  )
}

export default Docs