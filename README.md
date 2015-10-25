# gdsync
Google Drive client providing advanced filtering sync capabilities

This hack was born out of frustration encountered while trying to use
Google Drive as a host file system for software development projects.
The standard Google client thrashes trying to sync local files across
multiple systems in the face of the thousands of files created,
modified and removed during a typical automated build/test/install
cycle.  This thrashing can cause Google Drive to crash, slow to a
crawl, create multiple copies of a file and other nasty side effects
and renders the ultimate goal of using mulitple computers (travel
laptop, work server, home or work development systems, etc.)
independently to seamlessly continue working on one or more
development projects.

As a side benefit, this program provides a long sought, open source
Google Drive sync solution to the Linux community.  While doing this,
it also has adopted a multi-user model: one will typically have
multiple Google accounts, each of which may have the need to use a
Google Drive sync capability.

The initial target is to support filtering (ignore) patterns for
common souce control systems, especially Git, Bazaar, Mercurial,
Subversion, Perforce and CVS.  Later, general filtering will be
supported.

The design chosen is to provide simple command line operations that
build over time to provide a more complete solution.  Later, GUI and
Web based clients can be developed and used to leverage more polished
platform specific instances.  The inital commands will be things like:

    gdsync addUser pajatopmr@gmail.command
    gdsync showLog family@pajato.com
    gdsync killServer
    gdsync moveBase /Users/paul/GDS

Lastly, the design builds on the Google Drive API, which necessitates
using the Google Developers Console for authentication and Google
Drive access.  From this it follows that an entity can use the source
code to build a private or public hosted solution.  It also follows
that certain secure files (e.g. 'client_secret.json') must be protected
from inadvertently being stored in this repo.
