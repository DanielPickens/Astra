#this is a template spec and actual spec will be generated
#debuginfo not supported with Go
%global debug_package %{nil}
%global _enable_debug_package 0
%global __os_install_post /usr/lib/rpm/brp-compress %{nil}
%global package_name astra
%global product_name astra
%global golang_version ${GOLANG_VERSION}
%global golang_version_nastrat ${GOLANG_VERSION_NastraT}
%global astra_version ${astra_VERSION}
%global astra_rpm_version ${astra_RPM_VERSION}
%global astra_release ${astra_RELEASE}
%global git_commit  ${GIT_COMMIT}
%global astra_cli_version v%{astra_version}
%global source_dir openshift-astra-%{astra_version}-%{astra_release}
%global source_tar %{source_dir}.tar.gz
%global gopath  %{_builddir}/gocode
%global _missing_build_ids_terminate_build 0

Name:           %{package_name}
Version:        %{astra_rpm_version}
Release:        %{astra_release}%{?dist}
Summary:        %{product_name} client astra CLI binary
License:        ASL 2.0
URL:            https://github\.com/danielpickens/astra/tree/%{astra_cli_version}

Source0:        %{source_tar}
BuildRequires:  gcc
BuildRequires:  golang >= %{golang_version}
Provides:       %{package_name} = %{astra_rpm_version}
Obsoletes:      %{package_name} <= %{astra_rpm_version}

%description
astra is a fast, iterative, and straightforward CLI tool for developers who write, build, and deploy applications on OpenShift.

%prep
%setup -q -n %{source_dir}

%build
export GITCOMMIT="%{git_commit}"
mkdir -p %{gopath}/src/github.com/daniel-pickens
ln -s "$(pwd)" %{gopath}/src/github\.com/danielpickens/astra
export GOPATH=%{gopath}
cd %{gopath}/src/github\.com/danielpickens/astra
go mod edit -go=%{golang_version}
%ifarch x86_64
# go test -race is not supported on all arches
GOFLAGS='-mod=vendor' make test
%endif
make prepare-release
echo "%{astra_version}" > dist/release/VERSION
unlink %{gopath}/src/github\.com/danielpickens/astra

%install
mkdir -p %{buildroot}/%{_bindir}
install -m 0755 dist/bin/linux-`go env GOARCH`/astra %{buildroot}%{_bindir}/astra
mkdir -p %{buildroot}%{_datadir}
install -d %{buildroot}%{_datadir}/%{name}-redistributable
install -p -m 755 dist/release/astra-linux-amd64 %{buildroot}%{_datadir}/%{name}-redistributable/astra-linux-amd64
install -p -m 755 dist/release/astra-linux-arm64 %{buildroot}%{_datadir}/%{name}-redistributable/astra-linux-arm64
install -p -m 755 dist/release/astra-linux-ppc64le %{buildroot}%{_datadir}/%{name}-redistributable/astra-linux-ppc64le
install -p -m 755 dist/release/astra-linux-s390x %{buildroot}%{_datadir}/%{name}-redistributable/astra-linux-s390x
install -p -m 755 dist/release/astra-darwin-amd64 %{buildroot}%{_datadir}/%{name}-redistributable/astra-darwin-amd64
install -p -m 755 dist/release/astra-darwin-arm64 %{buildroot}%{_datadir}/%{name}-redistributable/astra-darwin-arm64
install -p -m 755 dist/release/astra-windows-amd64.exe %{buildroot}%{_datadir}/%{name}-redistributable/astra-windows-amd64.exe
cp -avrf dist/release/astra*.tar.gz %{buildroot}%{_datadir}/%{name}-redistributable
cp -avrf dist/release/astra*.zip %{buildroot}%{_datadir}/%{name}-redistributable
cp -avrf dist/release/SHA256_SUM %{buildroot}%{_datadir}/%{name}-redistributable
cp -avrf dist/release/VERSION %{buildroot}%{_datadir}/%{name}-redistributable

%files
%license LICENSE
%{_bindir}/astra

%package redistributable
Summary:        %{product_name} client CLI binaries for Linux, macOS and Windows
BuildRequires:  gcc
BuildRequires:  golang >= %{golang_version}
Provides:       %{package_name}-redistributable = %{astra_rpm_version}
Obsoletes:      %{package_name}-redistributable <= %{astra_rpm_version}

%description redistributable
%{product_name} client astra cross platform binaries for Linux, macOS and Windows.

%files redistributable
%license LICENSE
%dir %{_datadir}/%{name}-redistributable
%{_datadir}/%{name}-redistributable/astra-linux-amd64
%{_datadir}/%{name}-redistributable/astra-linux-amd64.tar.gz
%{_datadir}/%{name}-redistributable/astra-linux-arm64
%{_datadir}/%{name}-redistributable/astra-linux-arm64.tar.gz
%{_datadir}/%{name}-redistributable/astra-linux-ppc64le
%{_datadir}/%{name}-redistributable/astra-linux-ppc64le.tar.gz
%{_datadir}/%{name}-redistributable/astra-linux-s390x
%{_datadir}/%{name}-redistributable/astra-linux-s390x.tar.gz
%{_datadir}/%{name}-redistributable/astra-darwin-amd64
%{_datadir}/%{name}-redistributable/astra-darwin-amd64.tar.gz
%{_datadir}/%{name}-redistributable/astra-darwin-arm64
%{_datadir}/%{name}-redistributable/astra-darwin-arm64.tar.gz
%{_datadir}/%{name}-redistributable/astra-windows-amd64.exe
%{_datadir}/%{name}-redistributable/astra-windows-amd64.exe.zip
%{_datadir}/%{name}-redistributable/SHA256_SUM
%{_datadir}/%{name}-redistributable/VERSION
