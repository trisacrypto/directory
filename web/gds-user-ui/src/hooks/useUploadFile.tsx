import { useToast } from '@chakra-ui/toast';
import { useState } from 'react';
import { useUpdateCertificateStep } from './useUpdateCertificateStep';
import { StepEnum } from 'types/enums';
import { validationSchema } from 'modules/dashboard/certificate/lib';
import { handleError } from 'utils/utils';
const useUploadFile = () => {
  const { updateCertificateStep } = useUpdateCertificateStep();
  const toast = useToast();
  const [isFileLoading, setIsFileLoading] = useState<boolean>(false);
  const handleFileUpload = (file: any) => {
    setIsFileLoading(true);
    const reader = new FileReader();
    reader.onload = async (ev: any) => {
      const data = JSON.parse(ev.target.result);
      try {
        const validationData = await validationSchema[0].validate(data);
        const payload = {
          step: StepEnum.ALL,
          form: validationData
        };
        updateCertificateStep(payload);
        toast({
          title: 'File uploaded',
          description: 'Your file has been uploaded successfully',
          status: 'success',
          duration: 5000,
          isClosable: true,
          position: 'top-right'
        });
      } catch (e: any) {
        if (e.name === 'ValidationError') {
          toast({
            title: 'Invalid file',
            description: e.message || 'your json file is invalid',
            status: 'error',
            duration: 5000,
            isClosable: true,
            position: 'top-right'
          });
          handleError(e, `[Invalid file], it's missing some required fields : ${e.message}`);
        }
      } finally {
        setIsFileLoading(false);
      }
    };
    reader.readAsText(file);
  };

  return { isFileLoading, handleFileUpload };
};

export default useUploadFile;
