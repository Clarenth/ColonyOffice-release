import React from 'react'
import { useParams } from 'react-router-dom'

// Mutations
import { useGetDocumentByIDMutation, useGetFilesByDocumentIDMutation } from '@/lib/react-query/queriesAndMutations';

import Loader from '@/components/shared/Loader';
import DocsCreateForm from '@/components/forms/DocsCreateForm';

const DocsEdit = () => {
  const { document_id } = useParams();
  const { data: document, isPending } = useGetDocumentByIDMutation(document_id || '');
  const { data: files } = useGetFilesByDocumentIDMutation(document_id || '');

  if(isPending) return <Loader />

  return (
    <div className='flex flex-1'>
      <div className='common-container'>
        <div className='max-w-5xl flex-start gap-3 justify-start w-full'>
          <img 
            src='/assets/plus-dark-svgrepo-com.svg'
            alt='add-post'
            width={36}
            height={36}
          />
          <h2 className='h3-bold md:h2-bold text-left w-full'></h2>
        </div>
        <DocsCreateForm action="update" document={document} files={files}/>
      </div>
    </div>
  )
}

export default DocsEdit