# This file is used to initialize a new GraphDB repository. 
# It is filled at runtime.

@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#>.
@prefix rep: <http://www.openrdf.org/config/repository#>.
@prefix sr: <http://www.openrdf.org/config/repository/sail#>.
@prefix sail: <http://www.openrdf.org/config/sail#>.
@prefix graphdb: <http://www.ontotext.com/config/graphdb#>.

[] a rep:Repository ;
    rep:repositoryID "{{ .RepositoryID }}" ;
    rdfs:label "{{ .Label }}" ;
    rep:repositoryImpl [
        rep:repositoryType "graphdb:SailRepository" ;
        sr:sailImpl [
            sail:sailType "graphdb:Sail" ;

            graphdb:base-URL "{{ .BaseURL }}" ;
            graphdb:defaultNS "" ;
            graphdb:entity-index-size "10000000" ;
            graphdb:entity-id-size  "32" ;
            graphdb:imports "" ;
            graphdb:repository-type "file-repository" ;
            graphdb:ruleset "empty" ;
            graphdb:storage-folder "storage" ;

            graphdb:enable-context-index "true" ;

            graphdb:enablePredicateList "true" ;

            graphdb:in-memory-literal-properties "true" ;
            graphdb:enable-literal-index "true" ;

            graphdb:check-for-inconsistencies "false" ;
            graphdb:disable-sameAs  "true" ;
            graphdb:query-timeout  "0" ;
            graphdb:query-limit-results  "0" ;
            graphdb:throw-QueryEvaluationException-on-timeout "false" ;
            graphdb:read-only "false" ;
        ]
    ].