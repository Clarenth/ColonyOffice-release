import DocsCard from '@/components/shared/DocsCard';
import Loader from '@/components/shared/Loader';
import { useGetRecentDocumentsMutation } from '@/lib/react-query/queriesAndMutations';
import React from 'react'

const Home = () => {
  // const isLoadingDocs = true;
  // const documents = null;
  const { data: documents, isPending: isLoadingDocs, isError: isErrorDocs } = useGetRecentDocumentsMutation();

  return (
    <div className='flex flex-1'>
      <div className='home-container'>
        <div className='home-posts'>
          <h2 className='h3-bold md:h2-bold text-left w-full'>Home feed</h2>
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
  )
}

export default Home