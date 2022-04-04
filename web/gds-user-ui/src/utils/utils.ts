export const findStepKey = (steps: any, key: number) =>
  steps.filter((step: any) => step.key === key);

export const isValidUuid = (str: string) => {
  // Regular expression to check if string is a valid UUID
  const regexExp =
    /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/gi;
  return regexExp.test(str);
};

export const getStepData = (steps: any, key: number): TStep | undefined => {
  const s = findStepKey(steps, key);
  if (s && s?.length === 1) {
    return s[0].data;
  }
  return undefined;
};

export const hasStepError = (steps: any): boolean => {
  const s = steps.filter((step: any) => step.status === 'error');
  return s.length > 0;
};
