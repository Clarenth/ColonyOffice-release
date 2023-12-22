import * as zod from "zod"


export const SignupValidation = zod.object({
  email: zod.string().email(),
  password: zod.string().min(16, 'password is too short').max(128),
  phone_number: zod.string().min(2).max(20),
  job_title: zod.string().min(2).max(50),
  office_address: zod.string().min(2).max(50),
  employee_identity_data: zod.object({
    first_name: zod.string(),
    middle_name: zod.string(),
    last_name: zod.string(),
    sex: zod.string(),
    gender: zod.string(),
    age: zod.string(),
    height: zod.string(),
    home_address: zod.string(),
    birthdate: zod.string(),
    birthplace: zod.string(),
  }),
  security_access_level: zod.string().min(2).max(50),
})


export const LoginValidation = zod.object({
  email: zod.string().email(),
  password: zod.string().min(16, 'password is too short').max(128),
})

export const DocsValidation = zod.object({
  title: zod.string(),
  description: zod.string(),
  language: zod.string(),
  security_access_level: zod.string(),
  files: zod.custom<File[]>(),
})