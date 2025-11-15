package constants

const HelpDocumentation string = `Commands:
  2web new <project_name>
  2web generate <generator> <name>
  2web template <template>
  2web install <package_name>

  2web serve [path]
  2web build [path]
  2web lint [path]
  2web format [path]
  2web test [path]

  2web database <sub_command>
  2web doctor <sub_command>
  2web cms <sub_command>

  2web generate <template_name>
    component
    directive
    service
    aspect
    interceptor
    page
    guard
    model
    enum
    interface
    migration

  2web template <template_name>
    server-side-rendering
    database
    load-balancer
    sitemap
    robots.txt
    security.txt
    llms.txt

  2web database <sub_command>
    init
    migrate

  2web cms <sub_command>
    add
    view
    sync
    remove

  2web doctor <sub_command>
    check
    check-dependencies
    install-dependencies
`
