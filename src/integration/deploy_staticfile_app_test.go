package integration_test

import (
	"integration/cutlass"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deploy a staticfile app", func() {
	var app *cutlass.App
	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	BeforeEach(func() {
		app = cutlass.New(filepath.Join(bpDir, "cf_spec", "fixtures", "staticfile_app"))
		app.SetEnv("BP_DEBUG", "1")
	})

	It("", func() {
		Expect(app.Push()).To(Succeed())
		Expect(app.InstanceStates()).To(Equal([]string{"RUNNING"}))

		Expect(app.Stdout.String()).To(ContainSubstring("Buildpack version "))
		Expect(app.Stdout.String()).To(ContainSubstring("HOOKS 1: BeforeCompile"))
		Expect(app.Stdout.String()).To(ContainSubstring("HOOKS 2: AfterCompile"))
		Expect(app.Stdout.String()).To(MatchRegexp("nginx -p .*/nginx -c .*/nginx/conf/nginx.conf"))

		Expect(app.GetBody("/")).To(ContainSubstring("This is an example app for Cloud Foundry that is only static HTML/JS/CSS assets."))

		_, headers, err := app.Get("/fixture.json", map[string]string{})
		Expect(err).To(BeNil())
		Expect(headers["Content-Type"]).To(Equal([]string{"application/json"}))

		_, headers, err = app.Get("/lots_of.js", map[string]string{"Accept-Encoding": "gzip"})
		Expect(err).To(BeNil())
		Expect(headers).To(HaveKeyWithValue("Content-Encoding", []string{"gzip"}))

		By("requesting a non-compressed version of a compressed file", func() {
			By("with a client that can handle receiving compressed content", func() {
				By("returns and handles the file", func() {
					url, err := app.GetUrl("/war_and_peace.txt")
					Expect(err).To(BeNil())
					command := exec.Command("curl", "-s", "--compressed", url)
					Expect(command.Output()).To(ContainSubstring("Leo Tolstoy"))
				})
			})

			By("with a client that cannot handle receiving compressed content", func() {
				By("returns and handles the file", func() {
					url, err := app.GetUrl("/war_and_peace.txt")
					Expect(err).To(BeNil())
					command := exec.Command("curl", "-s", url)
					Expect(command.Output()).To(ContainSubstring("Leo Tolstoy"))
				})
			})
		})
	})

	PContext("with a cached buildpack", func() {
		// TODO :cached do
		It("logs the files it downloads", func() {
			// expect(app).to have_logged(/Copy \[\/.*\]/)
		})

		It("does not call out over the internet", func() {
			// expect(app).to_not have_internet_traffic
		})
	})

	PContext("with a uncached buildpack", func() {
		// TODO :uncached do
		It("logs the files it downloads", func() {
			// expect(app).to have_logged(/Download \[https:\/\/.*\]/)
		})

		It("uses a proxy during staging if present", func() {
			// expect(app).to use_proxy_during_staging
		})
	})

	Context("unpackaged buildpack eg. from github", func() {
		// let(:buildpack) { "staticfile-unpackaged-buildpack-#{rand(1000)}" }
		// let(:app) { Machete.deploy_app('staticfile_app', buildpack: buildpack, skip_verify_version: true) }
		// before do
		//   buildpack_file = "/tmp/#{buildpack}.zip"
		//   Open3.capture2e('zip','-r',buildpack_file,'bin/','src/', 'scripts/', 'manifest.yml','VERSION')[1].success? or raise 'Could not create unpackaged buildpack zip file'
		//   Open3.capture2e('cf', 'create-buildpack', buildpack, buildpack_file, '100', '--enable')[1].success? or raise 'Could not upload buildpack'
		//   FileUtils.rm buildpack_file
		// end
		// after do
		//   Open3.capture2e('cf', 'delete-buildpack', '-f', buildpack)
		// end

		PIt("runs", func() {
			// expect(app).to be_running
			// expect(app).to have_logged(/Running go build supply/)
			// expect(app).to have_logged(/Running go build finalize/)

			// browser.visit_path('/')
			// expect(browser).to have_body('This is an example app for Cloud Foundry that is only static HTML/JS/CSS assets.')
		})
	})

	Context("running a task", func() {
		BeforeEach(func() {
			// TODO
			// skip_if_no_run_task_support_on_targeted_cf
		})

		PIt("exits", func() {
			// expect(app).to be_running

			// Open3.capture2e('cf','run-task','staticfile_app','wc -l public/index.html')[1].success? or raise 'Could not create run task'
			// wait_until(60) do
			//   stdout, _ = Open3.capture2e('cf','tasks','staticfile_app')
			//   stdout =~ /SUCCEEDED.*wc.*index.html/
			// end
			// stdout, _ = Open3.capture2e('cf','tasks','staticfile_app')
			// expect(stdout).to match(/SUCCEEDED.*wc.*index.html/)
		})
	})
})
