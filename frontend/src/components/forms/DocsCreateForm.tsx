// Libraries
"use client"

import * as z from "zod" 
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { useNavigate } from "react-router-dom"

// Validation
import { DocsValidation } from "@/lib/validation"

// Mutation
import { useCreateDocumentMutation, useUpdateDocumentMutation, useUploadFileMutation } from "@/lib/react-query/queriesAndMutations"

// Components
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Textarea } from "../ui/textarea"
import FileUploader from "../shared/FileUploader"
import SecurityLevelDropdown from "../shared/SecurityLevelDropdown"
import { useToast } from "../ui/use-toast"

type DocsFormProps = {
  document?: {
    document_id: string;
    document_title: string;
    author_name: string;
    author_id: string;
    description: string;
    language: string;
    created_at: string;
    security_access_level: string;
  }

  files?: {
    file_title: string;
    author_name: string;
    security_access_level: string;
    created_at: string;
    updated_at: string;
  }

  action: 'create' | 'update';
}

const DocsCreateForm = ({ document, files, action }: DocsFormProps) => {
  const { mutateAsync: createDocument, isPending: isLoadingCreate } = useCreateDocumentMutation()
  const { mutateAsync: updateDocument, isPending: isLoadingUpdate} = useUpdateDocumentMutation();
  const { mutateAsync: uploadFile, isPending: isLoadingUpload } = useUploadFileMutation();
  
  const { toast } = useToast();
  const navigate = useNavigate()

  const form = useForm<z.infer<typeof DocsValidation>>({
    resolver: zodResolver(DocsValidation),
    defaultValues: {
      // This used to be a solution: title: document ? document?.document_title
      // but now it is not needed. Preserve it incase future me needs it.
      title: document ? document?.document.document_title : "",
      description: document ? document?.document.description : "",
      language: document ? document?.document.language : "",
      security_access_level: document ? document?.document.security_access_level : "",
      files: [],
    },
  })
  console.log("Hello from DocsCreateForm: ", document)

  // 2. Define a submit handler.
  async function onSubmit(values: z.infer<typeof DocsValidation>) {
    if(document && action === 'update') {
      const updatedDocument = await updateDocument({
        ...values,
        document_id: document.document_id,
      })
      if(!updatedDocument) {
        toast({
          title: "Updated failed. Try again."
        })
      }
      return navigate(`/docs/${document.document_id}`)
    }
    
    const formData = new FormData()
    const newDocument = await createDocument(values)

    if(!newDocument) {
      toast({
        title: "Create new document failed. Please try again."
      })
    }

    formData.append("document_id", newDocument)
    console.log(newDocument);

    for(let i = 0; i <= values.files.length; i++) {
      formData.append("files", values.files[i])
    }

    const newFiles = await uploadFile(formData)
    
    if(!newFiles) {
      toast({
        title: "Files upload failed. Investigate this error yourself.",
      })
    }

    // in future, once create is successful, navigate to the newly created document
    navigate("/doc/:id")
    //navigate("/")
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-9 w-full max-w-5xl">
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="shad-form_label">Title</FormLabel>
              <FormControl>
                <Input
                  type="text" 
                  placeholder="Title" 
                  className="shad-input"
                  {...field} 
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="shad-form_label">Description</FormLabel>
              <FormControl>
                <Textarea placeholder="shadcn" className="shad-textarea custom-scrollbar" {...field} />
              </FormControl>
              <FormMessage className="shad-form_message" />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="language"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="shad-form_label">Language</FormLabel>
              <FormControl>
                <Input
                  type="text"
                  placeholder="language" 
                  className="shad-input"
                  {...field} 
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="security_access_level"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="shad-form_label">Security Access Level</FormLabel>
              <FormControl>
                <SecurityLevelDropdown />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="files"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="shad-form_label">Upload Files</FormLabel>
              <FormControl>
                <FileUploader 
                  fieldChange={field.onChange}
                  mediaURL={ document?.fileURL }
                />
              </FormControl>
              <FormMessage className="shad-form_message" />
            </FormItem>
          )}
        />
        <div className="flex gap-4 items-center justify-end">
          <Button type="button" className="shad-button_dark_4">Cancel</Button>
          <Button 
            type="submit" 
            disabled={isLoadingCreate || isLoadingUpload} 
            className="shad-button_primary whitespace-nowrap"
          >
            { isLoadingCreate || isLoadingUpdate && 'Loading...' }
            Submit
          </Button>
        </div>
      </form>
    </Form>
  )
}

export default DocsCreateForm