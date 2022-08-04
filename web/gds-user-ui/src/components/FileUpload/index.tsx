import { ReactNode, useRef } from 'react';
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Icon,
  InputGroup,
  useColorModeValue
} from '@chakra-ui/react';
import { useForm, UseFormRegisterReturn } from 'react-hook-form';
import { FiFile } from 'react-icons/fi';
const MAX_FILE_SIZE = 10;
type FileUploadProps = {
  register: UseFormRegisterReturn;
  accept?: string;
  multiple?: boolean;
  children?: ReactNode;
  onFileSubmit: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onReaderLoader: (data: string) => void;
};
type FileUploaderProps = {
  onReadFileUploaded: (func: any) => void;
};
const FileUpload = (props: FileUploadProps) => {
  const { register, accept, multiple, children, onFileSubmit, onReaderLoader } = props;
  const inputRef = useRef<HTMLInputElement | null>(null);
  const { ref, ...rest } = register as { ref: (instance: HTMLInputElement | null) => void };

  const handleClick = () => inputRef.current?.click();

  return (
    <InputGroup onClick={handleClick}>
      <input
        type={'file'}
        multiple={multiple || false}
        hidden
        accept={accept}
        {...rest}
        ref={(e) => {
          ref(e);
          inputRef.current = e;
        }}
        onChange={(e) => onFileSubmit(e)}
      />
      <>{children}</>
    </InputGroup>
  );
};

type FormValues = {
  file_: FileList;
};

const FileUploader = ({ onReadFileUploaded }: FileUploaderProps) => {
  const {
    register,
    formState: { errors }
  } = useForm<FormValues>();
  const onSubmit = (e: any) => {
    e.preventDefault();
    console.log('submitted');
    const file = e.target.files[0];
    onReadFileUploaded(file);
  };

  const validateFiles = (value: FileList) => {
    if (value.length < 1) {
      return 'Files is required';
    }
    for (const file of Array.from(value)) {
      const fsMb = file.size / (1024 * 1024);
      if (fsMb > MAX_FILE_SIZE) {
        return 'Max file size 10mb';
      }
    }
    return true;
  };

  return (
    <>
      <form>
        <FormControl isInvalid={!!errors.file_} isRequired>
          <FileUpload
            accept={'.csv,application/JSON'}
            multiple
            onFileSubmit={onSubmit}
            onReaderLoader={onReadFileUploaded}
            register={register('file_', { validate: validateFiles })}>
            <Button
              bg={useColorModeValue('black', 'white')}
              color={useColorModeValue('white', 'black')}
              leftIcon={<Icon as={FiFile} />}
              minWidth={150}>
              Import File
            </Button>
          </FileUpload>

          <FormErrorMessage>{errors.file_ && errors?.file_.message}</FormErrorMessage>
        </FormControl>
      </form>
    </>
  );
};

export default FileUploader;
