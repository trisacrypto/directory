import { act, getByRole, render, screen } from "utils/test-utils";
import TrisaImplementationForm from "../TrisaImplementationForm";
import { dynamicActivate } from "utils/i18nLoaderHelper";
import TrisaImplementation from "..";

function renderComponent() {
    return render(<TrisaImplementationForm />);
}

jest.mock("hooks/useFetchCertificateStep", () => ({
    useFetchCertificateStep: () => ({
        certificateStep: {
            form: jest.fn(),
            errors: jest.fn(),
        },
        isFetchingCertificateStep: false,
    }),
}));

describe("TrisaImplementation", () => {
    beforeAll(() => {
        act(() => {
            dynamicActivate("en");
        });
    });

    it("should render", () => {
        const { container } = renderComponent();
        expect(container).toMatchSnapshot();
    });

    it("should render the form", () => {
        const { getByTestId } = renderComponent();
        const trisaImplementation = getByTestId("trisa-implementation-form");
        expect(trisaImplementation).toBeInTheDocument();
    });

    it("should render the TestNet header", () => {
        const { getByText } = renderComponent();
        const testNet = getByText("TRISA Endpoint: TestNet");
        expect(testNet).toBeInTheDocument();
    });

    it("should render the MainNet header", () => {
        const { getByText } = renderComponent();
        const mainNet = getByText("TRISA Endpoint: MainNet");
        expect(mainNet).toBeInTheDocument();
    });

    it("should render StepButtons", () => {
        const { getByTestId } = renderComponent();
        const stepButtons = getByTestId("step-buttons");
        expect(stepButtons).toBeInTheDocument();
    });
});