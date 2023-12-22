// Libraries
import React, { useCallback, useState } from 'react'
import { FileWithPath, useDropzone } from 'react-dropzone';

// Components
import { Button } from '../ui/button';

type FileUploaderProps = {
  fieldChange: (FILES: File[]) => void;
  mediaURL: string,
}

/*
TODO
-Replace the large image screen with a list of files coming from the files
 object we get with the documentID.
 Display: file_title
*/

const FileUploader = ({ fieldChange, mediaURL }: FileUploaderProps) => {
  const [file, setFile] = useState<File[]>([]);
  const [fileURL, setFileURL] = useState('');

  const onDrop = useCallback((acceptedFiles: FileWithPath[]) => {
    // Do something
    setFile(acceptedFiles);
    fieldChange(acceptedFiles);
    setFileURL(URL.createObjectURL(acceptedFiles[0]))
  }, [file])

  const { getRootProps, getInputProps } = useDropzone({
    onDrop,
    accept: {
      "image/*": ['.jpg', '.jpeg', '.png'],
      "application/txt": ['.txt',],
      "application/pdf": ['.pdf',],
      "application/odt": ['.odt',],
      "application/odp": ['.odp',],
      "application/ods": ['.ods',],
      "application/doc": ['.doc',],
      "application/docx": ['.docx',],
      "application/xls": ['.xls',],
      "application/xlsx": ['.xlsx',],
    }
  });

  return (
    <div {...getRootProps()} className='flex flex-center flex-col bg-dark-3 round-xl cursor-pointer'>
      <input {...getInputProps()} className='cursor-pointer'/>
      {
        fileURL ? 
        (
          <React.Fragment>
            <div className='flex flex-1 justify-center w-full p-5 lg:p-10'>
              <img 
                src={fileURL}
                alt='file to upload'
                //alt={}
                className='file_uploader-img'
              />
            </div>
            <p className='file_uploader-label'>Click or drag to upload</p>
          </React.Fragment>
        ) : 
        (
          <div className='file_uploader-box'>
            <img
              src='/assets/files-dark-svgrepo-com.svg'
              alt='file-upload'
              width={96}
              height={77}
            />
            <h3 className='base-medium text-light-2 mb-2 mt-6'>Drag photo here</h3>
            <p className='text-light-4 small-regular mb-6'>File Types</p>
            <Button className='shad-button_dark-4'>
              Select Files 
            </Button>
          </div>
        )
      }
    </div>
  )
}

export default FileUploader