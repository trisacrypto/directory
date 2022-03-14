import { Meta, Story } from "@storybook/react";
import { FormikProps } from "formik";
import BasicDetails from "./";

type BasicDetailsProps = {
  formik: FormikProps<any>;
};

export default {
  title: "components/BasicDetails",
  component: BasicDetails,
} as Meta<BasicDetailsProps>;

const Template: Story<BasicDetailsProps> = (args) => <BasicDetails {...args} />;

export const Default = Template.bind({});
Default.args = {};
