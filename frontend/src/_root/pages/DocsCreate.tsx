import CreateDocument from '@/components/forms/DocsCreateForm'
import React from 'react'

const DocsCreate = () => {
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
        <CreateDocument action='create' />
      </div>
    </div>
  )
}

export default DocsCreate